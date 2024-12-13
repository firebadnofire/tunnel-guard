# Build and installation variables
BINARY = tunnel-guard
SRC = main.go
SCRIPTS = tg-add-ssh-key tg-transfer-key
PREFIX ?= /usr
BINDIR = $(PREFIX)/bin
SYSTEMD_DIR = /etc/systemd/system
SERVICE_FILE = tunnel-guard.service
CONFIG_DIR = /etc/tunnel-guard

.PHONY: all build clean help install uninstall reinstall update \
	 install-binary install-scripts install-service \
	 uninstall-binary uninstall-scripts uninstall-service

all: help

help:
	@echo "Available targets:"
	@echo "  build	      - Build the $(BINARY) binary"
	@echo "  install	    - Install the $(BINARY), scripts, and systemd service"
	@echo "  uninstall	  - Uninstall the $(BINARY), scripts, systemd service, configs, keys, etc."
	@echo "  reinstall	  - Uninstall and then install (preserving user data)"
	@echo "  update	     - Alias for 'reinstall' (preserves user data)"
	@echo "  clean	      - Remove built binaries"

build:
	@echo "Building $(BINARY)..."
	@go build -o $(BINARY) $(SRC)

install: install-binary install-scripts install-service

install-binary: build
	@echo "Installing $(BINARY) to $(BINDIR)"
	@sudo install -d $(BINDIR)
	@sudo install -m 755 $(BINARY) $(BINDIR)/$(BINARY)
	@echo "Generating configuration by running $(BINARY) briefly..."
	@sudo sh -c '$(BINDIR)/$(BINARY) & sleep 2; kill $$! || true'

install-scripts:
	@echo "Installing scripts to $(BINDIR)"
	@for script in $(SCRIPTS); do \
		sudo install -m 755 $$script $(BINDIR)/$$script; \
	done

install-service:
	@echo "Installing systemd service file to $(SYSTEMD_DIR)"
	@sudo install -m 644 $(SERVICE_FILE) $(SYSTEMD_DIR)/$(SERVICE_FILE)
	@echo "Reloading systemd daemon..."
	@sudo systemctl daemon-reload
	@echo "Enabling and starting $(SERVICE_FILE)..."
	@sudo systemctl enable $(SERVICE_FILE)
	@sudo systemctl start $(SERVICE_FILE)

uninstall: uninstall-service uninstall-binary uninstall-scripts
	@echo "Removing configuration directory $(CONFIG_DIR)"
	@sudo rm -rf $(CONFIG_DIR)
	@echo "Removing user 'ssh-tun'"
	@sudo userdel -r ssh-tun || true
	@echo "Cleaning up built binaries"
	@rm -f $(BINARY)

uninstall-binary:
	@echo "Removing $(BINARY) from $(BINDIR)"
	@sudo rm -f $(BINDIR)/$(BINARY)

uninstall-scripts:
	@echo "Removing scripts from $(BINDIR)"
	@for script in $(SCRIPTS); do \
		sudo rm -f $(BINDIR)/$$script; \
	done

uninstall-service:
	@echo "Stopping and disabling $(SERVICE_FILE)..."
	@sudo systemctl stop $(SERVICE_FILE)
	@sudo systemctl disable $(SERVICE_FILE)
	@echo "Removing systemd service file from $(SYSTEMD_DIR)"
	@sudo rm -f $(SYSTEMD_DIR)/$(SERVICE_FILE)
	@echo "Reloading systemd daemon..."
	@sudo systemctl daemon-reload

# Preserve user data on reinstall
reinstall:
	@echo "Reinstalling without deleting user data..."
	@$(MAKE) uninstall-binary uninstall-scripts uninstall-service --no-print-directory
	@$(MAKE) install --no-print-directory

update: reinstall

clean:
	@echo "Cleaning up built binaries..."
	@rm -f $(BINARY)
