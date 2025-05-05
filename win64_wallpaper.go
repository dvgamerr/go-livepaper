package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

// WallpaperStyle defines the type for wallpaper display styles.
type WallpaperStyle int

const (
	STYLE_SPAN WallpaperStyle = 22
	STYLE_FILL WallpaperStyle = 10
)

const (
	spiSetDeskWallpaper = 0x0014
	spifUpdateINIFile   = 0x01
	spifSendChange      = 0x02
)

var (
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

func setWallpaper(imagePath string) error {
	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file not found: %s", imagePath)
	}

	imagePathPtr, err := syscall.UTF16PtrFromString(imagePath)
	if err != nil {
		return fmt.Errorf("failed to convert path to UTF16 pointer: %w", err)
	}

	// uiParam - not used for setting wallpaper path
	ret, _, err := systemParametersInfo.Call(spiSetDeskWallpaper, 0, uintptr(unsafe.Pointer(imagePathPtr)), spifUpdateINIFile|spifSendChange)

	if ret == 0 {
		// If ret is 0, check the error. If err is nil, it might still be an issue.
		// If err is not nil and not ERROR_SUCCESS, it's definitely an error.
		if err != nil && err.Error() != "The operation completed successfully." {
			return fmt.Errorf("failed to set wallpaper (SystemParametersInfo call failed): %w", err)
		}
		// It's possible ret is 0 but the operation succeeded (less common).
		// We might need more robust error checking depending on observed behavior.
		fmt.Println("SystemParametersInfo returned 0, but no specific error reported. Assuming success but monitor results.")
	} else {
		// Non-zero return usually means success. Check err just in case, though it's often nil here.
		if err != nil && err.Error() != "The operation completed successfully." {
			fmt.Printf("SystemParametersInfo returned non-zero (%d) but reported an error: %v. Proceeding cautiously.\n", ret, err)
		}
	}

	return nil
}

// setWallpaperStyle sets the desktop wallpaper style in the Windows registry.
func setWallpaperStyle(style WallpaperStyle) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()
	const tileValue = "0"

	if err = key.SetStringValue("WallpaperStyle", fmt.Sprintf("%d", style)); err != nil {
		return fmt.Errorf("failed to set WallpaperStyle registry value: %w", err)
	}

	if err = key.SetStringValue("TileWallpaper", tileValue); err != nil {
		return fmt.Errorf("failed to set TileWallpaper registry value: %w", err)
	}

	// Need to broadcast the change again after registry modification
	if ret, _, err := systemParametersInfo.Call(spiSetDeskWallpaper, 0, 0, spifUpdateINIFile|spifSendChange); ret == 0 {
		if err != nil && err.Error() != "The operation completed successfully." {
			return fmt.Errorf("failed to broadcast registry change (SystemParametersInfo call failed): %w", err)
		}
	} else {
		if err != nil && err.Error() != "The operation completed successfully." {
			fmt.Printf("SystemParametersInfo returned non-zero (%d) after style change but reported an error: %v. Proceeding cautiously.\n", ret, err)
		}
	}

	return nil
}
