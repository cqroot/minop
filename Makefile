.DEFAULT_GOAL := run

BUILD_DIR = $(CURDIR)/build

.PHONY: build
build:
	cmake -B '$(BUILD_DIR)' -S '$(CURDIR)' -G Ninja -DCMAKE_BUILD_TYPE=Debug -DCMAKE_EXPORT_COMPILE_COMMANDS=1
	ninja -C '$(BUILD_DIR)'
	@[ -f '$(CURDIR)/compile_commands.json' ] || ln -s '$(BUILD_DIR)/compile_commands.json' '$(CURDIR)/compile_commands.json'

.PHONY: run
run: build
	'$(BUILD_DIR)/minop'
