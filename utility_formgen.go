package fyneextensions

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
FormGenUtility is a struct that allows you to automatically generate
a form based on the fields of any given struct. It utilizes the Fyne library
to create these forms, efficiently creating a user interface for data input.

This utility is beneficial for rapidly constructing forms,
providing flexibility and consistent structuring to your applications.

Currently supported types are:
  - string
  - int, int8, int16, int32, int64
  - float32 and float64
  - bool
  - time.Time and time.Duration
  - slice of string []string
  - int alias for dropdown select entry only, e.g. type weekday int

Mainly, the Fyne form inputs are configured via the struct fields' tags.
Here are the tags available for customization:

  - `formGenInclude`: This tag determines whether a particular struct field
    should be included in the form. The value should be "true" if you want to
    include the field from the form. Skip it otherwise. only exported fields
    are allowed

  - `formGenDescription`: This tag provides a descriptive label to a form field.
    This will be shown side of the entry on the form

  - `formGenLabel`: This tag provides a descriptive label to a form field.
    This will be shown side of the check or radio items

  - `formGenMaxVal` and `formGenMinVal`: These tags determine the maximum and minimum
    values for numeric form fields. If out of range, it causes a validation error.

  - `formGenIsRequired`: This tag determines if a field is required. If a required
    field is left empty, it causes a validation error.

  - `formGenDefaultValue`: This tag defines a default value for a form field.
    For required fields ('formGenIsRequired=true'), this value will be used if entry is
    left empty.

  - `formGenOnEntryChange`: This tag allows defining a callback function that is triggered
    when the value of this form field changes. The function must be exported with struct
    receiver (not pointer receiver!) and must have as argument a value of same type of the
    field and an error. For float64, as example, callback should be:
    func (myStruct MyStruct) OnFloatChange(float64, error)
    Notice that there is no check at compile about validity of callback function: errors
    (e.g. spelling of the function) will result in panic

  - `formGenOptions`: This tag is used to define dropdown options for dropdown form fields.
    Options are applied only on custom struct derived from int, e.g. type MySelection int
    The MySelection values can be as defined with a const, e.g.
    const (
    Selection1 MySelection = iota,
    Selection2,
    ...
    )
    Options as should be shown in the dropdown select entry are defined as options in this tag.
    Each option should be separated by "|||"
    for example `formGenOptions` "Selection 1|||Selection 2|||Selection 3"

  - `formGenRadioGroup`, 'formGenCheckGroup': This tag allows defining a group of check
    or radio buttons (on a single line). The tag value should be a unique string which
    defines the group.

Each field of the struct corresponds to a field on the form,
and the types, labels, and other behaviors of the form fields are
mapped from the types, names, and tags of the struct fields.

# The method ShowDialog will generate a dialog.Form and show it on the passed canvas

Them method GetWidget will generate a new fyne.Widget to be shown in a custom CanvasObject
*/
type FormGenUtility struct {
	s interface{}

	w          fyne.Window
	firstEntry fyne.Focusable

	formItems     []*widget.FormItem
	fieldSetter   map[string]func()
	formClearer   map[string]func()
	entrySetter   map[string]func()
	entryDisabler map[string]func()
	entryEnabler  map[string]func()

	formDialog *dialog.FormDialog
	OnSubmit   func(bool)
	formWidget *widget.Form

	firstShow bool
}

// OverrideEntry overrides the content of the entry from the value stored in the specified field.
// Field name must mach the field as typed in the struct passed to the FormGenUtility. Use to
// update the entry in the form if the struct is altered, while the form is open
func (fgu *FormGenUtility) OverrideEntry(fieldName string) {
	if f, ok := fgu.entrySetter[fieldName]; ok {
		f()
	}
}

// OverrideField overrides the value stored in the struct from the value in the dedicated entry.
// Field name must mach the field as typed in the struct passed to the FormGenUtility. Use to
// update the struct from the entry before the Confirm button is pressed in the form
func (fgu *FormGenUtility) OverrideField(fieldName string) {
	if f, ok := fgu.fieldSetter[fieldName]; ok {
		f()
	}
}

// EnableEntry enables the entry in the form.
// Field name must mach the field as typed in the struct passed to the FormGenUtility.
func (fgu *FormGenUtility) EnableEntry(fieldName string) {
	if f, ok := fgu.entryEnabler[fieldName]; ok {
		f()
	}
}

// DisableEntry disables the entry in the form
// Field name must mach the field as typed in the struct passed to the FormGenUtility.
func (fgu *FormGenUtility) DisableEntry(fieldName string) {
	if f, ok := fgu.entryDisabler[fieldName]; ok {
		f()
	}
}

func (fgu *FormGenUtility) createFormItems() {
	v := reflect.Indirect(reflect.ValueOf(fgu.s))
	t := v.Type()

	firstEntrySet := false
	radioGroupEntries := make(map[string]*widget.RadioGroup)
	radioGroupFields := make(map[string]map[string]reflect.Value)
	checkGroupEntries := make(map[string]*widget.CheckGroup)
	checkGroupFields := make(map[string]map[string]reflect.Value)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() {
			fieldName := field.Name

			fieldValue := v.Field(i)
			fieldSetterFunc := func() {}
			fieldClearerFunc := func() {}
			entrySetterFunc := func() {}
			entryDisabler := func() {}
			entryEnabler := func() {}

			//fmt.Println(fieldValue.Type().String())

			include, includeOk := field.Tag.Lookup("formGenInclude")
			if includeOk && strings.ToLower(include) != "false" {
				description := field.Tag.Get("formGenDescription")
				if description == "" {
					description = fieldName
				}
				label := field.Tag.Get("formGenLabel")
				if label == "" {
					label = fieldName
				}

				hiLimStr := field.Tag.Get("formGenMaxVal")
				loLimStr := field.Tag.Get("formGenMinVal")

				requiredField := false
				if v, ok := field.Tag.Lookup("formGenIsRequired"); ok && strings.ToLower(v) != "false" {
					requiredField = true
				}
				defaultValStrg := field.Tag.Get("formGenDefaultValue")
				isPassword := false
				if v, ok := field.Tag.Lookup("formGenIsPassword"); ok && strings.ToLower(v) != "false" {
					isPassword = true
				}

				onChangeCallbackStrg := field.Tag.Get("formGenOnEntryChange")
				onChangeCallbackFunc := v.MethodByName(onChangeCallbackStrg)

				dropDown, dropDownOk := field.Tag.Lookup("formGenOptions")

				radioGroupName, radioGroupOk := field.Tag.Lookup("formGenRadioGroup")
				checkGroupName, checkGroupOk := field.Tag.Lookup("formGenCheckGroup")

				var entry fyne.CanvasObject
				floatBS := 0
				intBS := -1

				switch fieldValue.Type() {
				case reflect.TypeOf(""):
					mEntry := widget.NewEntry()
					if !firstEntrySet {
						fgu.firstEntry = mEntry
						firstEntrySet = true
					}

					if isPassword {
						mEntry.Password = true
					}

					defaultVal := defaultValStrg
					mEntry.SetPlaceHolder(defaultVal)

					textConverter := func(s string) (string, error) {
						if s == "" {
							if requiredField {
								return "", fmt.Errorf("empty required field")
							}
							return defaultVal, nil
						}
						return s, nil
					}

					mEntry.Validator = func(s string) error {
						_, err := textConverter(s)
						return err
					}

					if onChangeCallbackFunc.IsValid() {
						mEntry.OnChanged = func(s string) {
							if v, err := textConverter(s); err == nil {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
							} else {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(err)})
							}
						}
					}

					fieldSetterFunc = func() {
						if v, err := textConverter(mEntry.Text); err == nil {
							fieldValue.SetString(v)
						}
					}
					fieldClearerFunc = func() {
						fieldValue.SetString("")
						mEntry.SetText("")
					}
					entrySetterFunc = func() {
						mEntry.SetText(fieldValue.String())
					}
					entryDisabler = mEntry.Disable
					entryEnabler = mEntry.Enable
					entry = mEntry
				case reflect.TypeOf(int(0)):
					intBS = 0
					fallthrough
				case reflect.TypeOf(int8(0)):
					if intBS == -1 {
						intBS = 8
					}
					fallthrough
				case reflect.TypeOf(int16(0)):
					if intBS == -1 {
						intBS = 16
					}
					fallthrough
				case reflect.TypeOf(int32(0)):
					if intBS == -1 {
						intBS = 32
					}
					fallthrough
				case reflect.TypeOf(int64(0)):
					if intBS == -1 {
						intBS = 64
					}

					mEntry := widget.NewEntry()
					if !firstEntrySet {
						fgu.firstEntry = mEntry
						firstEntrySet = true
					}

					var defaultVal, minVal, maxVal int64
					maxVal = math.MaxInt64
					minVal = math.MinInt64

					if v, err := strconv.ParseInt(defaultValStrg, 0, intBS); err == nil {
						defaultVal = v
						mEntry.SetPlaceHolder(defaultValStrg)
					}
					if v, err := strconv.ParseInt(loLimStr, 0, intBS); err == nil {
						minVal = v
					}
					if v, err := strconv.ParseInt(hiLimStr, 0, intBS); err == nil {
						maxVal = v
					}

					textConverter := func(s string) (int, int8, int16, int32, int64, error) {
						if s == "" {
							if requiredField {
								return 0, 0, 0, 0, 0, fmt.Errorf("empty required field")
							}
							return int(defaultVal), int8(defaultVal), int16(defaultVal), int32(defaultVal), defaultVal, nil
						}
						if v, err := strconv.ParseInt(s, 10, intBS); err == nil {
							if v > maxVal {
								return int(v), int8(v), int16(v), int32(v), v, fmt.Errorf("out of range, max: %s", hiLimStr)
							}
							if v < minVal {
								return int(v), int8(v), int16(v), int32(v), v, fmt.Errorf("out of range, min: %s", loLimStr)
							}
							return int(v), int8(v), int16(v), int32(v), v, nil
						}
						return 0, 0, 0, 0, 0, fmt.Errorf("invalid duration : %s", s)
					}

					mEntry.Validator = func(s string) error {
						_, _, _, _, _, err := textConverter(s)
						return err
					}

					if onChangeCallbackFunc.IsValid() {
						mEntry.OnChanged = func(s string) {
							if v, v8, v16, v32, v64, err := textConverter(s); err == nil {
								switch intBS {
								case 0:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
								case 8:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v8), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
								case 16:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v16), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
								case 32:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v32), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
								case 64:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v64), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
								}
							} else {
								switch intBS {
								case 0:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(err)})
								case 8:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v8), reflect.ValueOf(err)})
								case 16:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v16), reflect.ValueOf(err)})
								case 32:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v32), reflect.ValueOf(err)})
								case 64:
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v64), reflect.ValueOf(err)})
								}
							}
						}
					}

					fieldSetterFunc = func() {
						if _, _, _, _, v64, err := textConverter(mEntry.Text); err == nil {
							fieldValue.SetInt(v64)
						}
					}
					fieldClearerFunc = func() {
						fieldValue.SetInt(0)
						mEntry.SetText("")
					}
					entrySetterFunc = func() {
						mEntry.SetText(strconv.FormatInt(fieldValue.Int(), 10))
					}
					entryDisabler = mEntry.Disable
					entryEnabler = mEntry.Enable
					entry = mEntry
				case reflect.TypeOf(float32(0.)):
					floatBS = 32
					fallthrough
				case reflect.TypeOf(float64(0.)):
					if floatBS == 0 {
						floatBS = 64
					}

					mEntry := widget.NewEntry()
					if !firstEntrySet {
						fgu.firstEntry = mEntry
						firstEntrySet = true
					}

					var defaultVal, minVal, maxVal float64
					if v, err := strconv.ParseFloat(defaultValStrg, floatBS); err == nil {
						defaultVal = v
						mEntry.SetPlaceHolder(defaultValStrg)
					}
					minVal = math.Inf(-1)
					if v, err := strconv.ParseFloat(loLimStr, floatBS); err == nil {
						minVal = v
					}
					maxVal = math.Inf(1)
					if v, err := strconv.ParseFloat(hiLimStr, floatBS); err == nil {
						maxVal = v
					}
					mEntry.SetPlaceHolder(defaultValStrg)

					textConverter := func(s string) (float64, float32, error) {
						if s == "" {
							if requiredField {
								return math.NaN(), float32(math.NaN()), fmt.Errorf("empty required field")
							}
							return defaultVal, float32(defaultVal), nil
						}
						if v, err := strconv.ParseFloat(s, floatBS); err == nil {
							if v > maxVal {
								return v, float32(v), fmt.Errorf("out of range, max: %g", maxVal)
							}
							if v < minVal {
								return v, float32(v), fmt.Errorf("out of range, min: %g", minVal)
							}
							return v, float32(v), nil
						}
						return math.NaN(), float32(math.NaN()), fmt.Errorf("not a number : %s", s)
					}

					mEntry.Validator = func(s string) error {
						_, _, err := textConverter(s)
						return err
					}

					if onChangeCallbackFunc.IsValid() {
						mEntry.OnChanged = func(s string) {
							if v64, v32, err := textConverter(s); err == nil {
								if floatBS == 32 {
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v32), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
								} else {
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v64), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
								}
							} else {
								if floatBS == 32 {
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v32), reflect.ValueOf(err)})
								} else {
									onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v64), reflect.ValueOf(err)})
								}
							}
						}
					}

					fieldSetterFunc = func() {
						if v64, _, err := textConverter(mEntry.Text); err == nil {
							fieldValue.SetFloat(v64)
						}
					}
					fieldClearerFunc = func() {
						fieldValue.SetFloat(0.)
						mEntry.SetText("")
					}
					entrySetterFunc = func() {
						mEntry.SetText(strconv.FormatFloat(fieldValue.Float(), 'g', -1, floatBS))
					}
					entryDisabler = mEntry.Disable
					entryEnabler = mEntry.Enable
					entry = mEntry
				case reflect.TypeOf(true):
					if checkGroupOk {
						var callbackFunc func([]string)
						if onChangeCallbackFunc.IsValid() {
							callbackFunc = func(b []string) {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(b), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
							}
						}

						var mEntry *widget.CheckGroup
						if mCheckGroupEntry, ok := checkGroupEntries[checkGroupName]; ok {
							mEntry = mCheckGroupEntry
							mEntry.OnChanged = callbackFunc
						} else {
							mEntry = widget.NewCheckGroup([]string{}, callbackFunc)
							mEntry.Horizontal = true
							checkGroupEntries[checkGroupName] = mEntry
							checkGroupFields[checkGroupName] = make(map[string]reflect.Value)
							entry = mEntry

						}
						mEntry.Append(label)
						checkGroupFields[checkGroupName][label] = fieldValue

						//if !firstEntrySet {
						//	fgu.firstEntry = mEntry
						//	firstEntrySet = true
						//}

						fieldSetterFunc = func() {
							ttp := mEntry.Selected
							checkGroupFields[checkGroupName][label].SetBool(false)
							for _, s := range ttp {
								if s == label {
									checkGroupFields[checkGroupName][label].SetBool(true)
								}
							}
						}
						fieldClearerFunc = func() {
							checkGroupFields[checkGroupName][label].SetBool(false)
							mEntry.SetSelected([]string{})
						}
						entrySetterFunc = func() {
							selected := make([]string, 0)
							for s, f := range checkGroupFields[checkGroupName] {
								if f.Bool() {
									selected = append(selected, s)
								}
							}
							mEntry.SetSelected(selected)
						}
						entryDisabler = mEntry.Disable
						entryEnabler = mEntry.Enable

						defaultValL := strings.ToLower(defaultValStrg)
						if defaultValL == "true" {
							mEntry.Selected = append(mEntry.Selected, label) //.SetChecked(true)
						}
					} else if radioGroupOk {
						var callbackFunc func(string)
						if onChangeCallbackFunc.IsValid() {
							callbackFunc = func(b string) {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(b), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
							}
						}

						var mEntry *widget.RadioGroup
						if mRadioGroupEntry, ok := radioGroupEntries[radioGroupName]; ok {
							mEntry = mRadioGroupEntry
							mEntry.OnChanged = callbackFunc
						} else {
							mEntry = widget.NewRadioGroup([]string{}, callbackFunc)
							mEntry.Horizontal = true
							radioGroupEntries[radioGroupName] = mEntry
							radioGroupFields[radioGroupName] = make(map[string]reflect.Value)
							entry = mEntry

						}
						mEntry.Append(label)
						radioGroupFields[radioGroupName][label] = fieldValue

						//if !firstEntrySet {
						//	fgu.firstEntry = mEntry
						//	firstEntrySet = true
						//}

						fieldSetterFunc = func() {
							ttp := mEntry.Selected
							if ttp == label {
								radioGroupFields[radioGroupName][label].SetBool(true)
							} else {
								radioGroupFields[radioGroupName][label].SetBool(false)
							}
						}
						fieldClearerFunc = func() {
							radioGroupFields[radioGroupName][label].SetBool(false)
							mEntry.SetSelected("")
						}
						entrySetterFunc = func() {
							selected := ""
							for s, f := range radioGroupFields[radioGroupName] {
								if f.Bool() {
									selected = s
									break
								}
							}
							mEntry.SetSelected(selected)
						}
						entryDisabler = mEntry.Disable
						entryEnabler = mEntry.Enable

						defaultValL := strings.ToLower(defaultValStrg)
						if defaultValL == "true" {
							mEntry.Selected = label //.SetChecked(true)
						}
					} else {
						var callbackFunc func(bool)
						if onChangeCallbackFunc.IsValid() {
							callbackFunc = func(b bool) {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(b), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
							}
						}

						mEntry := widget.NewCheck(label, callbackFunc)
						if !firstEntrySet {
							fgu.firstEntry = mEntry
							firstEntrySet = true
						}

						defaultValL := strings.ToLower(defaultValStrg)
						if defaultValL == "true" {
							mEntry.SetChecked(true)
						}
						fieldSetterFunc = func() {
							ttp := mEntry.Checked
							fieldValue.SetBool(ttp)
						}
						fieldClearerFunc = func() {
							fieldValue.SetBool(false)
							mEntry.SetChecked(false)
						}
						entrySetterFunc = func() {
							mEntry.SetChecked(fieldValue.Bool())
						}
						entryDisabler = mEntry.Disable
						entryEnabler = mEntry.Enable
						entry = mEntry
					}
				case reflect.TypeOf(time.Time{}):
					mEntry := widget.NewEntry()
					if !firstEntrySet {
						fgu.firstEntry = mEntry
						firstEntrySet = true
					}

					var defaultVal time.Time
					if strings.ToLower(defaultValStrg) == "now" || defaultValStrg == "" {
						defaultVal = time.Now()
					} else {
						if v, err := time.Parse("02-Jan-2006 15:04:05 MST", defaultValStrg); err == nil {
							defaultVal = v
						} else {
							if v, err := time.ParseInLocation("02-Jan-2006 15:04:05", defaultValStrg, time.Local); err == nil {
								defaultVal = v
							}
						}
					}
					if !defaultVal.IsZero() {
						mEntry.SetPlaceHolder(defaultVal.Format("02-Jan-2006 15:04:05 MST"))
					}

					var maxVal time.Time
					if v, err := time.Parse("02-Jan-2006 15:04:05 MST", hiLimStr); err == nil {
						maxVal = v
					} else {
						if v, err := time.ParseInLocation("02-Jan-2006 15:04:05", hiLimStr, time.Local); err == nil {
							maxVal = v
						}
					}

					var minVal time.Time
					if v, err := time.Parse("02-Jan-2006 15:04:05 MST", loLimStr); err == nil {
						minVal = v
					} else {
						if v, err := time.ParseInLocation("02-Jan-2006 15:04:05", loLimStr, time.Local); err == nil {
							minVal = v
						}
					}

					textConverter := func(s string) (time.Time, error) {
						if s == "" {
							if requiredField {
								return time.Time{}, fmt.Errorf("empty required field")
							}
							return defaultVal, nil
						}

						if v, err := time.Parse("02-Jan-2006 15:04:05 MST", s); err == nil {
							if !maxVal.IsZero() {
								if v.After(maxVal) {
									return v, fmt.Errorf("out of range, max: %s", hiLimStr)
								}
							}
							if !minVal.IsZero() {
								if v.Before(minVal) {
									return v, fmt.Errorf("out of range, min: %s", loLimStr)
								}
							}
							return v, nil
						}
						if v, err := time.ParseInLocation("02-Jan-2006 15:04:05", s, time.Local); err == nil {
							if !maxVal.IsZero() {
								if v.After(maxVal) {
									return v, fmt.Errorf("out of range, max: %s", hiLimStr)
								}
							}
							if !minVal.IsZero() {
								if v.Before(minVal) {
									return v, fmt.Errorf("out of range, min: %s", loLimStr)
								}
							}
							return v, nil
						}
						return time.Time{}, fmt.Errorf("not a timestamp : %s", s)
					}

					mEntry.Validator = func(s string) error {
						_, err := textConverter(s)
						return err
					}

					if onChangeCallbackFunc.IsValid() {
						mEntry.OnChanged = func(s string) {
							if v, err := textConverter(s); err == nil {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
							} else {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(err)})
							}
						}
					}

					fieldSetterFunc = func() {
						if v, err := textConverter(mEntry.Text); err == nil {
							fieldValue.Set(reflect.ValueOf(v))
						}
					}
					fieldClearerFunc = func() {
						v := time.Time{}
						fieldValue.Set(reflect.ValueOf(v))
						mEntry.SetText("")
					}
					entrySetterFunc = func() {
						if v, ok := fieldValue.Interface().(time.Time); ok {
							if !v.IsZero() {
								mEntry.SetText(v.Format("02-Jan-2006 15:04:05 MST"))
							} else {
								mEntry.SetText("")
							}
						}
					}
					entryDisabler = mEntry.Disable
					entryEnabler = mEntry.Enable
					entry = mEntry
				case reflect.TypeOf(time.Duration(0)):
					mEntry := widget.NewEntry()
					if !firstEntrySet {
						fgu.firstEntry = mEntry
						firstEntrySet = true
					}

					var defaultVal, minVal, maxVal time.Duration
					minVal = math.MinInt64
					maxVal = math.MaxInt64

					if v, err := ParseDurationExtended(defaultValStrg); err == nil {
						defaultVal = v
						mEntry.SetPlaceHolder(defaultValStrg)
					}
					if v, err := ParseDurationExtended(loLimStr); err == nil {
						minVal = v
					}
					if v, err := ParseDurationExtended(hiLimStr); err == nil {
						maxVal = v
					}

					textConverter := func(s string) (time.Duration, error) {
						if s == "" {
							if requiredField {
								return time.Duration(0), fmt.Errorf("empty required field")
							}
							return defaultVal, nil
						}

						if v, err := ParseDurationExtended(s); err == nil {
							if v > maxVal {
								return v, fmt.Errorf("out of range, max: %s", hiLimStr)
							}
							if v < minVal {
								return v, fmt.Errorf("out of range, min: %s", loLimStr)
							}
							return v, nil
						}
						return time.Duration(0), fmt.Errorf("invalid duration : %s", s)
					}

					mEntry.Validator = func(s string) error {
						_, err := textConverter(s)
						return err
					}

					if onChangeCallbackFunc.IsValid() {
						mEntry.OnChanged = func(s string) {
							if v, err := textConverter(s); err == nil {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
							} else {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(err)})
							}
						}
					}

					fieldSetterFunc = func() {
						if v, err := textConverter(mEntry.Text); err == nil {
							fieldValue.Set(reflect.ValueOf(v))
						}
					}
					fieldClearerFunc = func() {
						v := time.Duration(0)
						fieldValue.Set(reflect.ValueOf(v))
						mEntry.SetText("")
					}
					entrySetterFunc = func() {
						if v, ok := fieldValue.Interface().(time.Duration); ok {
							if v.Nanoseconds() != 0 {
								mEntry.SetText(formatDuration(v))
							} else {
								mEntry.SetText("")
							}
						}
					}
					entryDisabler = mEntry.Disable
					entryEnabler = mEntry.Enable

					entry = mEntry
				case reflect.TypeOf([]string{}):
					mEntry := widget.NewMultiLineEntry()
					if !firstEntrySet {
						fgu.firstEntry = mEntry
						firstEntrySet = true
					}

					defaultVal := make([]string, 0)
					if defaultValStrg != "" {
						if v := strings.Split(defaultValStrg, "\n"); len(v) > 0 {
							defaultVal = v
							mEntry.SetPlaceHolder(defaultValStrg)
						}
					}

					textConverter := func(s string) ([]string, error) {
						if s == "" {
							if requiredField {
								return []string{}, fmt.Errorf("empty required field")
							}
							return defaultVal, nil
						}
						retArr := strings.Split(s, "\n")
						for i, s := range retArr {
							retArr[i] = strings.Trim(s, " \r\n\t")
						}
						if retArr[len(retArr)-1] == "" {
							retArr = retArr[:len(retArr)-1]
						}
						return retArr, nil
					}

					mEntry.Validator = func(s string) error {
						_, err := textConverter(s)
						return err
					}

					if onChangeCallbackFunc.IsValid() {
						mEntry.OnChanged = func(s string) {
							if v, err := textConverter(s); err == nil {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
							} else {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(err)})
							}
						}
					}

					fieldSetterFunc = func() {
						if v, err := textConverter(mEntry.Text); err == nil {
							fieldValue.Set(reflect.ValueOf(v))
						}
					}
					fieldClearerFunc = func() {
						v := make([]string, 0)
						fieldValue.Set(reflect.ValueOf(v))
						mEntry.SetText("")
					}
					entrySetterFunc = func() {
						fs := ""
						if v, ok := fieldValue.Interface().([]string); ok {
							for _, s := range v {
								fs = fs + s + "\n"
							}
						}
						mEntry.SetText(fs)
					}
					entryDisabler = mEntry.Disable
					entryEnabler = mEntry.Enable
					entry = mEntry
				default:
					if fieldValue.Kind() == reflect.TypeOf(int(0)).Kind() && dropDownOk {
						options := strings.Split(dropDown, ";")

						var callbackFunc func(string)
						if onChangeCallbackFunc.IsValid() {
							callbackFunc = func(b string) {
								onChangeCallbackFunc.Call([]reflect.Value{reflect.ValueOf(b), reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())})
							}
						}
						mEntry := widget.NewSelect(options, callbackFunc)
						if !firstEntrySet {
							fgu.firstEntry = mEntry
							firstEntrySet = true
						}

						var defaultVal int

						if v, err := strconv.ParseInt(defaultValStrg, 0, 0); err == nil {
							if v >= 0 && int(v) < len(options) {
								defaultVal = int(v)
								mEntry.SetSelectedIndex(int(v))
							}
						} else {
							defaultVal = 0
							mEntry.PlaceHolder = "Default: " + options[0]
						}

						textConverter := func(id int) (int, error) {
							if id < 0 {
								if requiredField {
									return 0, fmt.Errorf("empty required field")
								}
								return defaultVal, nil
							}
							if id >= len(options) {
								if requiredField {
									return 0, fmt.Errorf("empty required field")
								}
								return defaultVal, nil
							}

							return id, nil
						}

						fieldSetterFunc = func() {
							if v, err := textConverter(mEntry.SelectedIndex()); err == nil {
								fieldValue.SetInt(int64(v))
							}
						}
						fieldClearerFunc = func() {
							fieldValue.SetInt(int64(0))
							mEntry.SetSelected("")
						}
						entrySetterFunc = func() {
							mEntry.SetSelectedIndex(int(fieldValue.Int()))
						}
						entryDisabler = mEntry.Disable
						entryEnabler = mEntry.Enable
						entry = mEntry
					} else {
						continue
					}
				}

				fgu.fieldSetter[fieldName] = fieldSetterFunc
				fgu.formClearer[fieldName] = fieldClearerFunc
				fgu.entrySetter[fieldName] = entrySetterFunc
				fgu.entryEnabler[fieldName] = entryEnabler
				fgu.entryDisabler[fieldName] = entryDisabler

				if entry != nil {
					nItem := widget.NewFormItem(description, entry)
					fgu.formItems = append(fgu.formItems, nItem)
				}
			}
		}
	}
}

/*
NewFormGenDialog is a utility function for dynamically generating form dialog. It takes in a struct and creates a form dialog based on the struct's exported fields that are flagged with a `formGenInclude` tag.

Parameters:

- s (interface{}): The struct based on which the form dialog will be created.
- title (string): The title of the form dialog.
- confirm (string): The label for the confirm button in the form dialog.
- dismiss (string): The label for the cancel/dismiss button in the form dialog.
- onSubmit (func(bool)): A callback function that will be called when form is submitted. The function takes a boolean parameter which is true if form submission was successful and false otherwise.
- w (fyne.Window): The parent window for the form dialog.
- minSize (fyne.Size): The minimum size of the form dialog.

Returns:

- *FormGenUtility: A pointer to the FormGenUtility instance, which includes the generated form dialog and its associated elements.

Example Usage:

	type User struct {
	    Name string `formGenInclude:"true"`
	    Age int `formGenInclude:"true"`
	}

	user := User{Name: "John", Age: 25}

	fgu := NewFormGenDialog(user, "User Form", "OK", "Cancel", func(b bool) {
	    if b {
	        fmt.Println("Form Submitted")
	    } else {
	        fmt.Println("Form Dismissed")
	    }
	}, window, fyne.NewSize(300, 200))

In this example, a new form dialog is generated based on the `User` struct. The form will have two fields, "Name" and "Age". Both fields are included because they have `formGenInclude` tag set to "true". Upon form submission or cancellation, respective messages will be printed on the console.

The `NewFormGenDialog` function allows for the dynamic creation of forms without having to manually create each form field, making code cleaner and facilitating use of forms in larger projects.
*/
func NewFormGenDialog(s interface{}, title, confirm, dismiss string, onSubmit func(bool), w fyne.Window, minSize fyne.Size) *FormGenUtility {
	fgu := &FormGenUtility{
		s:             s,
		w:             w,
		formItems:     make([]*widget.FormItem, 0),
		fieldSetter:   make(map[string]func()),
		formClearer:   make(map[string]func()),
		entrySetter:   make(map[string]func()),
		entryEnabler:  make(map[string]func()),
		entryDisabler: make(map[string]func()),
		OnSubmit:      onSubmit,
		firstShow:     true,
	}

	fgu.createFormItems()

	fgu.formDialog = dialog.NewForm(title, confirm, dismiss, fgu.formItems, func(b bool) {
		if b {
			for _, o := range fgu.fieldSetter {
				o()
			}
		}
		if onSubmit != nil {
			onSubmit(b)
		}
	}, w)

	fgu.formWidget = widget.NewForm(fgu.formItems...)
	fgu.formWidget.SubmitText = confirm
	fgu.formWidget.CancelText = dismiss
	fgu.formWidget.OnSubmit = func() {
		for _, o := range fgu.fieldSetter {
			o()
		}
		if fgu.OnSubmit != nil {
			fgu.OnSubmit(true)
		}
	}
	fgu.formWidget.OnCancel = func() {
		if fgu.OnSubmit != nil {
			fgu.OnSubmit(false)
		}
	}

	mw, mh := fgu.formDialog.MinSize().Width, fgu.formDialog.MinSize().Height
	if minSize.Width > mw {
		mw = minSize.Width
	}
	if minSize.Height > mh {
		mh = minSize.Height
	}
	fgu.formDialog.Resize(fyne.NewSize(mw, mh))

	return fgu
}

type FormGenKeepValueOption int

const (
	KeepStructExceptOnFirstShow FormGenKeepValueOption = iota
	KeepStruct
	KeepEntry
	ResetStruct
)

/*
GetDialog is a method of the FormGenUtility type that returns the generated form dialog.

Parameters:
  - keepStructValues (FormGenKeepValueOption): A flag to indicate whether to keep or discard the values in the form entries.

Returns:

- *dialog.FormDialog: The form dialog generated by the NewFormGenDialog function.

Example Usage:

	fgu := NewFormGenDialog(user, "User Form", "OK", "Cancel", onSubmitFunc, window, fyne.NewSize(300, 200))

	// Save current values and get the form dialog
	dialog := fgu.GetDialog(true)
*/
func (fgu *FormGenUtility) GetDialog(keepStructValues FormGenKeepValueOption) *dialog.FormDialog {
	switch keepStructValues {
	case KeepStructExceptOnFirstShow:
		if !fgu.firstShow {
			for _, o := range fgu.entrySetter {
				o()
			}
		}
	case KeepStruct:
		for _, o := range fgu.entrySetter {
			o()
		}
	case KeepEntry:

	case ResetStruct:
		for _, o := range fgu.formClearer {
			o()
		}
	}

	fgu.firstShow = false
	return fgu.formDialog
}

/*
ShowDialog is a method of the FormGenUtility type that shows the generated form dialog.

Parameters:
  - keepStructValues (FormGenKeepValueOption): A flag to indicate whether to keep or discard the values in the form entries.

Example Usage:

	fgu := NewFormGenDialog(user, "User Form", "OK", "Cancel", onSubmitFunc, window, fyne.NewSize(300, 200))

	// Show dialog with current values
	fgu.ShowDialog(true)
*/
func (fgu *FormGenUtility) ShowDialog(keepStructValues FormGenKeepValueOption) {
	switch keepStructValues {
	case KeepStructExceptOnFirstShow:
		if !fgu.firstShow {
			for _, o := range fgu.entrySetter {
				o()
			}
		}
	case KeepStruct:
		for _, o := range fgu.entrySetter {
			o()
		}
	case KeepEntry:

	case ResetStruct:
		for _, o := range fgu.formClearer {
			o()
		}
	}

	fgu.firstShow = false

	fgu.formDialog.Show()
	if fgu.firstEntry != nil {
		fgu.w.Canvas().Focus(fgu.firstEntry)
	}
}

/*
GetWidget is a method of the FormGenUtility type that returns the generated form widget.
Differently from dialog.Form (GetDialog), the widget can be added to a custom CanvasObject
and further customized where necessary. It does not include ok and cancel buttons, which
need to be created manually and call the OnSubmit (or OnCancel) functions

Parameters:
  - keepStructValues (FormGenKeepValueOption): A flag to indicate whether to keep or discard the values in the form entries.

Returns:
- *widget.Form: The form widget generated by the NewFormGenDialog function.

Example Usage:

	fgu := NewFormGenDialog(user, "User Form", "OK", "Cancel", onSubmitFunc, window, fyne.NewSize(300, 200))

	// Retrieve the form widget with current values
	widget := fgu.GetWidget(true)
*/
func (fgu *FormGenUtility) GetWidget(keepStructValues FormGenKeepValueOption) *widget.Form {
	switch keepStructValues {
	case KeepStructExceptOnFirstShow:
		if !fgu.firstShow {
			for _, o := range fgu.entrySetter {
				o()
			}
		}
	case KeepStruct:
		for _, o := range fgu.entrySetter {
			o()
		}
	case KeepEntry:

	case ResetStruct:
		for _, o := range fgu.formClearer {
			o()
		}
	}

	fgu.firstShow = false

	return fgu.formWidget
}

/*
FocusOnFirstEntry is a method of the FormGenUtility type that performs the action of focussing on the first entry field of the form.

Example Usage:

	fgu := NewFormGenDialog(user, "User Form", "OK", "Cancel", onSubmitFunc, window, fyne.NewSize(300, 200))

	// Focus on the first entry field of the form
	fgu.FocusOnFirstEntry()

The FocusOnFirstEntry method does not have any parameters and does not return any value.

This method checks if there is at least one entry in the form, selects the first and then
sets the focus to this field. This action is especially useful after the form dialog is shown
to the user, to immediately allow the user to start entering data into the form without
having to manually click on the first entry field.
*/
func (fgu *FormGenUtility) FocusOnFirstEntry() {
	if fgu.firstEntry != nil {
		fgu.w.Canvas().Focus(fgu.firstEntry)
	}
}

var errLeadingInt = fmt.Errorf("time: bad [0-9]*") // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt[bytes []byte | string](s bytes) (x uint64, rem bytes, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, rem, errLeadingInt
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, rem, errLeadingInt
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x uint64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + uint64(c) - '0'
		if y > 1<<63 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}

var unitMap = map[string]uint64{
	"ns": uint64(time.Nanosecond),
	"us": uint64(time.Microsecond),
	"µs": uint64(time.Microsecond), // U+00B5 = micro symbol
	"μs": uint64(time.Microsecond), // U+03BC = Greek letter mu
	"ms": uint64(time.Millisecond),
	"s":  uint64(time.Second),
	"m":  uint64(time.Minute),
	"h":  uint64(time.Hour),
	"d":  uint64(time.Hour) * 24,
	"w":  uint64(time.Hour) * 24 * 7,
	"mo": uint64(time.Hour) * 24 * 30,
	"y":  uint64(time.Hour) * 24 * 365,
}

// ParseDurationExtended parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDurationExtended(s string) (time.Duration, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d uint64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, fmt.Errorf("time: invalid duration: %s", orig)
	}
	for s != "" {
		var (
			v, f  uint64      // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return 0, fmt.Errorf("time: invalid duration: %s", orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, fmt.Errorf("time: invalid duration: %s", orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, fmt.Errorf("time: invalid duration %s", orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, fmt.Errorf("time: missing unit in duration: %s", orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, fmt.Errorf("time: unknown unit %s in duration : %s", u, orig)
		}
		if v > 1<<63/unit {
			// overflow
			return 0, fmt.Errorf("time: invalid duration :%s", orig)
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += uint64(float64(f) * (float64(unit) / scale))
			if v > 1<<63 {
				// overflow
				return 0, fmt.Errorf("time: invalid duration %s", orig)
			}
		}
		d += v
		if d > 1<<63 {
			return 0, fmt.Errorf("time: invalid duration: %s", orig)
		}
	}
	if neg {
		return -time.Duration(d), nil
	}
	if d > 1<<63-1 {
		return 0, fmt.Errorf("time: invalid duration, %s", orig)
	}
	return time.Duration(d), nil
}

// format formats the representation of d into the end of buf and
// returns the offset of the first character.
func formatDuration(d time.Duration) string {
	// Largest time is 2540400h10m10.000000000s
	w := 200
	buf := make([]byte, w)

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	if u < uint64(time.Second) {
		// Special case: if duration is smaller than a second,
		// use smaller units, like 1.2ms
		var prec int
		w--
		buf[w] = 's'
		w--
		switch {
		case u == 0:
			buf[w] = '0'
			return string(buf[w:])
		case u < uint64(time.Microsecond):
			// print nanoseconds
			prec = 0
			buf[w] = 'n'
		case u < uint64(time.Millisecond):
			// print microseconds
			prec = 3
			// U+00B5 'µ' micro sign == 0xC2 0xB5
			w-- // Need room for two bytes.
			copy(buf[w:], "µ")
		default:
			// print milliseconds
			prec = 6
			buf[w] = 'm'
		}
		w, u = fmtFrac(buf[:w], u, prec)
		w = fmtInt(buf[:w], u)
	} else {
		w--
		buf[w] = 's'

		w, u = fmtFrac(buf[:w], u, 9)

		// u is now integer seconds
		w = fmtInt(buf[:w], u%60)
		u /= 60

		// u is now integer minutes
		if u > 0 {
			w--
			buf[w] = 'm'
			w = fmtInt(buf[:w], u%60)
			u /= 60

			// u is now integer hours
			if u > 0 {
				w--
				buf[w] = 'h'
				w = fmtInt(buf[:w], u%24)
				u /= 24

				// u is now integer days
				if u > 0 {
					w--
					buf[w] = 'd'
					w = fmtInt(buf[:w], u%365)
					u /= 365

					// u is now integer years
					if u > 0 {
						w--
						buf[w] = 'y'
						w = fmtInt(buf[:w], u)
					}
				}
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros. It omits the decimal
// point too when the fraction is 0. It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}
