/*
 * MinIO Cloud Storage, (C) 2019 MinIO, Inc.
 * PGG Obstor, (C) 2021-2026 PGG, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package color

import (
	"fmt"
)

// global colors.
var (
	// Check if we stderr, stdout are dumb terminals, we do not apply
	// ansi coloring on dumb terminals.
	IsTerminal = func() bool {
		return !NoColor
	}

	Bold = func() func(a ...interface{}) string {
		if IsTerminal() {
			return New(AttrBold).SprintFunc()
		}
		return fmt.Sprint
	}()

	RedBold = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(AttrFgRed, AttrBold).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	Red = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(AttrFgRed).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	Blue = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(FgBlue).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	Yellow = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(FgYellow).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	Green = func() func(a ...interface{}) string {
		if IsTerminal() {
			return New(FgGreen).SprintFunc()
		}
		return fmt.Sprint
	}()

	GreenBold = func() func(a ...interface{}) string {
		if IsTerminal() {
			return New(FgGreen, AttrBold).SprintFunc()
		}
		return fmt.Sprint
	}()

	CyanBold = func() func(a ...interface{}) string {
		if IsTerminal() {
			return New(FgCyan, AttrBold).SprintFunc()
		}
		return fmt.Sprint
	}()

	YellowBold = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(FgYellow, AttrBold).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	BlueBold = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(FgBlue, AttrBold).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	BgYellow = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(AttrBgYellow).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	Black = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(FgBlack).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	FgRed = func() func(a ...interface{}) string {
		if IsTerminal() {
			return New(AttrFgRed).SprintFunc()
		}
		return fmt.Sprint
	}()

	BgRed = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(AttrBgRed).SprintfFunc()
		}
		return fmt.Sprintf
	}()

	FgWhite = func() func(format string, a ...interface{}) string {
		if IsTerminal() {
			return New(AttrFgWhite).SprintfFunc()
		}
		return fmt.Sprintf
	}()
)
