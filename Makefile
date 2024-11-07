# Build and installation variables
BINARY = tunnel-guard
SRC = main.go
PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
SYSTEMD_DIR = /etc/systemd/system
SERVICE_FILE = tunnel-guard.service
CONFIG_DIR = /etc/tunnel-guard

.PHONY: all build clean help install uninstall reinstall update

all: help

help:
	@echo "Available targets:"
	@echo "  build         - Build the $(BINARY) binary"
	@echo "  install       - Install the $(BINARY) binary and systemd service"
	@echo "  uninstall     - Uninstall the $(BINARY) binary, systemd service, configs, keys, etc."
	@echo "  reinstall     - Uninstall and then install the $(BINARY)"
	@echo "  update        - Alias for 'reinstall'"
	@echo "  clean         - Remove built binaries"

build:
	@echo "Building $(BINARY)..."
	@go build -o $(BINARY) $(SRC)

install: build
	@echo "Installing $(BINARY) to $(BINDIR)"
	@sudo install -d $(BINDIR)
	@sudo install -m 755 $(BINARY) $(BINDIR)/$(BINARY)
	
	@echo "Generating configuration by running $(BINARY) briefly..."
	@sudo sh -c '$(BINDIR)/$(BINARY) & sleep 2; kill $$! || true'
	
	@echo "Installing systemd service file to $(SYSTEMD_DIR)"
	@sudo install -m 644 $(SERVICE_FILE) $(SYSTEMD_DIR)/$(SERVICE_FILE)
	@echo "Reloading systemd daemon..."
	@sudo systemctl daemon-reload
	@echo "Enabling and starting $(SERVICE_FILE)..."
	@sudo systemctl enable $(SERVICE_FILE)
	@sudo systemctl start $(SERVICE_FILE)

uninstall:
	@echo "Stopping and disabling $(SERVICE_FILE)..."
	@sudo systemctl stop $(SERVICE_FILE)
	@sudo systemctl disable $(SERVICE_FILE)
	@echo "Removing $(BINARY) from $(BINDIR)"
	@sudo rm -f $(BINDIR)/$(BINARY)
	@echo "Removing systemd service file from $(SYSTEMD_DIR)"
	@sudo rm -f $(SYSTEMD_DIR)/$(SERVICE_FILE)
	@echo "Reloading systemd daemon..."
	@sudo systemctl daemon-reload
	@echo "Removing configuration directory $(CONFIG_DIR)"
	@sudo rm -rf $(CONFIG_DIR)
	@echo "Removing user 'ssh-tun'"
	@sudo userdel -r ssh-tun || true
	@echo "Cleaning up built binaries"
	@rm -f $(BINARY)

reinstall: uninstall install

update: reinstall

clean:
	@echo "Cleaning up built binaries..."
	@rm -f $(BINARY)
