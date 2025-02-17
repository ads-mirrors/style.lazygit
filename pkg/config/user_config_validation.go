package config

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/constants"
)

func (config *UserConfig) Validate() error {
	if err := validateEnum("gui.statusPanelView", config.Gui.StatusPanelView,
		[]string{"dashboard", "allBranchesLog"}); err != nil {
		return err
	}
	if err := validateEnum("gui.showDivergenceFromBaseBranch", config.Gui.ShowDivergenceFromBaseBranch,
		[]string{"none", "onlyArrow", "arrowAndNumber"}); err != nil {
		return err
	}
	if err := validateKeybindings(config.Keybinding); err != nil {
		return err
	}
	return nil
}

func validateEnum(name string, value string, allowedValues []string) error {
	if slices.Contains(allowedValues, value) {
		return nil
	}
	allowedValuesStr := strings.Join(allowedValues, ", ")
	return fmt.Errorf("Unexpected value '%s' for '%s'. Allowed values: %s", value, name, allowedValuesStr)
}

func validateKeybindingsRecurse(path string, node any) error {
	value := reflect.ValueOf(node)
	if value.Kind() == reflect.Struct {
		for _, field := range reflect.VisibleFields(reflect.TypeOf(node)) {
			var newPath string
			if len(path) == 0 {
				newPath = field.Name
			} else {
				newPath = fmt.Sprintf("%s.%s", path, field.Name)
			}
			if err := validateKeybindingsRecurse(newPath,
				value.FieldByName(field.Name).Interface()); err != nil {
				return err
			}
		}
	} else if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			if err := validateKeybindingsRecurse(
				fmt.Sprintf("%s[%d]", path, i), value.Index(i).Interface()); err != nil {
				return err
			}
		}
	} else if value.Kind() == reflect.String {
		key := node.(string)
		if !isValidKeybindingKey(key) {
			return fmt.Errorf("Unrecognized key '%s' for keybinding '%s'. For permitted values see %s",
				strings.ToLower(key), path, constants.Links.Docs.CustomKeybindings)
		}
	} else {
		panic("Unexpected type")
	}
	return nil
}

func validateKeybindings(keybindingConfig KeybindingConfig) error {
	return validateKeybindingsRecurse("", keybindingConfig)
}
