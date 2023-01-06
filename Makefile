BUILD_DIR = ./release
TEMP_DIR = ../temp/battery-history.exe
VERSION=$(shell date /t)
build:
	cd cmd && go build -o $(TEMP_DIR)
	go-msi make --msi ./release/battery-history.msi --version $(VERSION) --src ./installer/windows/
	del .\temp\battery-history.exe
	rmdir .\temp