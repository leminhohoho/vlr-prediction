// NOTE: package borrowed from https://github.com/elulcao/progress-bar
//
//	This a modified version that curated for this specific project
package progressbar

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
)

const (
	// Reset
	Reset = "\x1B[0m"

	// Regular Colors
	Black  = "\x1B[30m"
	Red    = "\x1B[31m"
	Green  = "\x1B[32m"
	Yellow = "\x1B[33m"
	Blue   = "\x1B[34m"
	Purple = "\x1B[35m"
	Cyan   = "\x1B[36m"
	White  = "\x1B[37m"

	// Bold Colors
	BoldBlack  = "\x1B[1;30m"
	BoldRed    = "\x1B[1;31m"
	BoldGreen  = "\x1B[1;32m"
	BoldYellow = "\x1B[1;33m"
	BoldBlue   = "\x1B[1;34m"
	BoldPurple = "\x1B[1;35m"
	BoldCyan   = "\x1B[1;36m"
	BoldWhite  = "\x1B[1;37m"

	// Background Colors
	BgBlack  = "\x1B[40m"
	BgRed    = "\x1B[41m"
	BgGreen  = "\x1B[42m"
	BgYellow = "\x1B[43m"
	BgBlue   = "\x1B[44m"
	BgPurple = "\x1B[45m"
	BgCyan   = "\x1B[46m"
	BgWhite  = "\x1B[47m"
)

// PBar is the progress bar model
type PBar struct {
	Total        uint64         // Total number of iterations to sum 100%
	headerLength uint16         // Header length, to be used to calculate the bar width "Progress: [100%] []"
	headerText   string         // The text that is displayed before the counter, if none is set "Progress" will be displayed instead
	wscol        uint16         // Window width
	wsrow        uint16         // Window height
	doneStr      string         // Progress bar done string
	ongoingStr   string         // Progress bar ongoing string
	signalWinch  chan os.Signal // Signal handler: SIGWINCH
	signalTerm   chan os.Signal // Signal handler: SIGTERM
	once         sync.Once      // Close the signal channel only once
	winSize      struct {       // winSize is the struct to store the current window size, used by ioctl
		Row    uint16 // row
		Col    uint16 // column
		Xpixel uint16 // X pixel
		Ypixel uint16 // Y pixel
	}
}

// NewPBar create a new progress bar
func NewPBar() *PBar {
	pb := &PBar{
		Total:        0,
		headerLength: 0,
		wscol:        0,
		wsrow:        0,
		doneStr:      "█",
		ongoingStr:   " ",
		signalWinch:  make(chan os.Signal, 1),
		signalTerm:   make(chan os.Signal, 1),
	}

	signal.Notify(pb.signalWinch, syscall.SIGWINCH) // Register SIGWINCH signal
	signal.Notify(pb.signalTerm, syscall.SIGTERM)   // Register SIGTERM signal

	_ = pb.UpdateWSize()

	return pb
}

// CleanUp restore reserved bottom line and restore cursor position
func (pb *PBar) CleanUp() {
	pb.once.Do(func() { close(pb.signalWinch) }) // Close the signal channel politely, avoid closing it twice
	pb.once.Do(func() { close(pb.signalTerm) })  // Close the signal channel politely

	if pb.winSize.Col == 0 || pb.winSize.Row == 0 {
		return // Not a terminal, running in a pipeline or test
	}

	fmt.Print("\x1B7")                 // Save the cursor position
	fmt.Printf("\x1B[0;%dr", pb.wsrow) // Drop margin reservation
	fmt.Printf("\x1B[%d;0f", pb.wsrow) // Move the cursor to the bottom line
	fmt.Print("\x1B[0K")               // Erase the entire line
	fmt.Print("\x1B8")                 // Restore the cursor position util new size is calculated
}

// updateWSize update the window size
func (pb *PBar) UpdateWSize() error {
	isTerminal, err := pb.checkIsTerminal()
	if err != nil {
		return fmt.Errorf("could not check if the current process is running in a terminal: %w", err)
	}
	if !isTerminal {
		return nil // Not a terminal, running in a pipeline or test
	}
	if pb.Total == uint64(100) {
		return nil // No need to update the header length
	}

	pb.wscol = pb.winSize.Col
	pb.wsrow = pb.winSize.Row

	var header string

	if pb.headerText == "" {
		header = "Progress"
	} else {
		header = pb.headerText
	}

	headerDefaultLength := 7 + len(header) + helpers.CountNumLength(int(pb.Total))*2

	switch {
	case pb.wscol >= uint16(0) && pb.wscol <= uint16(20):
		pb.headerLength = uint16(0) // len("[100%]") is the minimum header length
	default:
		pb.headerLength = uint16(
			headerDefaultLength,
		) // len("Progress: [100%] []") is the maximum header length
	}

	fmt.Print("\x1BD")                   // Return carriage
	fmt.Print("\x1B7")                   // Save the cursor position
	fmt.Printf("\x1B[0;%dr", pb.wsrow-1) // Reserve the bottom line
	fmt.Print("\x1B8")                   // Restore the cursor position
	fmt.Print("\x1B[1A")                 // Moves cursor up # lines

	return nil
}

// Set a new header text which is displayed before the counter
func (pb *PBar) SetHeaderText(headerText string) {
	pb.headerText = headerText
	_ = pb.UpdateWSize()
}

// Set the total count of the progress bar
func (pb *PBar) SetTotalCount(total int) {
	pb.Total = uint64(total)
	_ = pb.UpdateWSize()
}

// Return the total count of the progress bar
func (pb *PBar) TotalCount() int {
	return int(pb.Total)
}

// RenderPBar render the progress bar. Receives the current iteration count
func (pb *PBar) RenderPBar(count int) {
	if pb.winSize.Col == 0 || pb.winSize.Row == 0 {
		return // Not a terminal, running in a pipeline or test
	}

	fmt.Print("\x1B7")       // Save the cursor position
	fmt.Print("\x1B[2K")     // Erase the entire line
	fmt.Print("\x1B[0J")     // Erase from cursor to end of screen
	fmt.Print("\x1B[?47h")   // Save screen
	fmt.Print("\x1B[1J")     // Erase from cursor to beginning of screen
	fmt.Print("\x1B[?47l")   // Restore screen
	defer fmt.Print("\x1B8") // Restore the cursor position util new size is calculated

	ratio := float64(count) / float64(pb.Total)

	barWidth := int(math.Abs(float64(pb.wscol - pb.headerLength))) // Calculate the bar width
	barDone := int(float64(barWidth) * ratio)                      // Calculate the bar done length
	done := strings.Repeat(pb.doneStr, barDone)                    // Fill the bar with done string
	todo := strings.Repeat(pb.ongoingStr, barWidth-barDone)        // Fill the bar with todo string
	bar := fmt.Sprintf("%s%s", done, todo)                         // Combine the done and todo string

	var highlight string

	if ratio < 0.3 {
		highlight = Red
	} else if ratio >= 0.3 && ratio < 0.7 {
		highlight = Yellow
	} else {
		highlight = Green
	}

	fmt.Printf("\x1B[%d;%dH", pb.wsrow, 0) // move cursor to row #, col #

	switch {
	case pb.wscol >= uint16(0) && pb.wscol <= uint16(9):
		fmt.Printf("[\x1B[33m%3d%%\x1B[0m]", uint64(count)*100/pb.Total)
	case pb.wscol >= uint16(10) && pb.wscol <= uint16(20):
		fmt.Printf("[\x1B[33m%3d%%\x1B[0m] %s", uint64(count)*100/pb.Total, bar)
	default:
		fmt.Printf("%s: %s%d/%d %s%s", pb.headerText, highlight, count, pb.Total, bar, Reset)
	}
}

// SignalHandler handle the signals, like SIGWINCH and SIGTERM
func (pb *PBar) SignalHandler() {
	go func() {
		for {
			select {
			case <-pb.signalWinch:
				if err := pb.UpdateWSize(); err != nil {
					panic(err) // The window size could not be updated
				}
			case <-pb.signalTerm:
				pb.CleanUp() // Restore reserved bottom line
				os.Exit(1)   // Exit gracefully but exit code 1
			}
		}
	}()
}

// checkIsTerminal check if the current process is running in a terminal
func (pb *PBar) checkIsTerminal() (isTerminal bool, err error) {
	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&pb.winSize))); err != 0 {
		if err == syscall.ENOTTY || err == syscall.ENODEV {
			return false, nil // Not a terminal, running in a pipeline or test
		} else {
			return false, err // Other error
		}
	}

	return true, nil
}
