// Copyright 2019 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cc

import (
	"android/soong/android"
	"testing"
)

func TestLinkerScript(t *testing.T) {
	t.Run("script", func(t *testing.T) {
		testCc(t, `
		cc_object {
			name: "foo",
			srcs: ["baz.o"],
			linker_script: "foo.lds",
		}`)
	})
}

func TestCcObjectWithBazel(t *testing.T) {
	bp := `
cc_object {
	name: "foo",
	srcs: ["baz.o"],
	bazel_module: { label: "//foo/bar:bar" },
}`
	config := TestConfig(t.TempDir(), android.Android, nil, bp, nil)
	config.BazelContext = android.MockBazelContext{
		OutputBaseDir: "outputbase",
		LabelToOutputFiles: map[string][]string{
			"//foo/bar:bar": []string{"bazel_out.o"}}}
	ctx := testCcWithConfig(t, config)

	module := ctx.ModuleForTests("foo", "android_arm_armv7-a-neon").Module()
	outputFiles, err := module.(android.OutputFileProducer).OutputFiles("")
	if err != nil {
		t.Errorf("Unexpected error getting cc_object outputfiles %s", err)
	}

	expectedOutputFiles := []string{"outputbase/execroot/__main__/bazel_out.o"}
	android.AssertDeepEquals(t, "output files", expectedOutputFiles, outputFiles.Strings())
}
