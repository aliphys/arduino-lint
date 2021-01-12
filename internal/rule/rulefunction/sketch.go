// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of Arduino Lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package rulefunction

// The rule functions for sketches.

import (
	"strings"

	"github.com/arduino/arduino-cli/arduino/globals"
	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/project/sketch"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
)

// SketchNameMismatch checks for mismatch between sketch folder name and primary file name.
func SketchNameMismatch() (result ruleresult.Type, output string) {
	for extension := range globals.MainFileValidExtensions {
		validPrimarySketchFilePath := projectdata.ProjectPath().Join(projectdata.ProjectPath().Base() + extension)
		exist, err := validPrimarySketchFilePath.ExistCheck()
		if err != nil {
			panic(err)
		}

		if exist {
			return ruleresult.Pass, ""
		}
	}

	return ruleresult.Fail, projectdata.ProjectPath().Base() + ".ino"
}

// ProhibitedCharactersInSketchFileName checks for prohibited characters in the sketch file names.
func ProhibitedCharactersInSketchFileName() (result ruleresult.Type, output string) {
	directoryListing, _ := projectdata.ProjectPath().ReadDir()
	directoryListing.FilterOutDirs()

	foundInvalidSketchFileNames := []string{}
	for _, potentialSketchFile := range directoryListing {
		if sketch.HasSupportedExtension(potentialSketchFile) {
			if !validProjectPathBaseName(potentialSketchFile.Base()) {
				foundInvalidSketchFileNames = append(foundInvalidSketchFileNames, potentialSketchFile.Base())
			}
		}
	}

	if len(foundInvalidSketchFileNames) > 0 {
		return ruleresult.Fail, strings.Join(foundInvalidSketchFileNames, ", ")
	}

	return ruleresult.Pass, ""
}

// SketchFileNameGTMaxLength checks if the sketch file names exceed the maximum length.
func SketchFileNameGTMaxLength() (result ruleresult.Type, output string) {
	directoryListing, _ := projectdata.ProjectPath().ReadDir()
	directoryListing.FilterOutDirs()

	foundTooLongSketchFileNames := []string{}
	for _, potentialSketchFile := range directoryListing {
		if sketch.HasSupportedExtension(potentialSketchFile) {
			if len(potentialSketchFile.Base())-len(potentialSketchFile.Ext()) > 63 {
				foundTooLongSketchFileNames = append(foundTooLongSketchFileNames, potentialSketchFile.Base())
			}
		}
	}

	if len(foundTooLongSketchFileNames) > 0 {
		return ruleresult.Fail, strings.Join(foundTooLongSketchFileNames, ", ")
	}

	return ruleresult.Pass, ""
}

// PdeSketchExtension checks for use of deprecated .pde sketch file extensions.
func PdeSketchExtension() (result ruleresult.Type, output string) {
	directoryListing, _ := projectdata.ProjectPath().ReadDir()
	directoryListing.FilterOutDirs()
	pdeSketches := []string{}
	for _, filePath := range directoryListing {
		if filePath.Ext() == ".pde" {
			pdeSketches = append(pdeSketches, filePath.Base())
		}
	}

	if len(pdeSketches) > 0 {
		return ruleresult.Fail, strings.Join(pdeSketches, ", ")
	}

	return ruleresult.Pass, ""
}

// IncorrectSketchSrcFolderNameCase checks for incorrect case of src subfolder name in recursive format libraries.
func IncorrectSketchSrcFolderNameCase() (result ruleresult.Type, output string) {
	directoryListing, err := projectdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "src")
	if found {
		return ruleresult.Fail, path.String()
	}

	return ruleresult.Pass, ""
}

// SketchDotJSONJSONFormat checks whether the sketch.json metadata file is a valid JSON document.
func SketchDotJSONJSONFormat() (result ruleresult.Type, output string) {
	metadataPath := sketch.MetadataPath(projectdata.ProjectPath())
	if metadataPath == nil {
		return ruleresult.Skip, "No metadata file"
	}

	if isValidJSON(metadataPath) {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, ""
}

// SketchDotJSONFormat checks whether the sketch.json metadata file has the required data format.
func SketchDotJSONFormat() (result ruleresult.Type, output string) {
	metadataPath := sketch.MetadataPath(projectdata.ProjectPath())
	if metadataPath == nil {
		return ruleresult.Skip, "No metadata file"
	}

	if projectdata.MetadataLoadError() == nil {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, projectdata.MetadataLoadError().Error()
}
