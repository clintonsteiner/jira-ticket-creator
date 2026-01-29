.PHONY: install serve build clean help

help:
	@echo "Documentation Build Commands"
	@echo "============================="
	@echo ""
	@echo "  make install    Install Jekyll dependencies"
	@echo "  make serve      Run docs locally (http://localhost:4000)"
	@echo "  make build      Build static site to _site/"
	@echo "  make clean      Remove built site"
	@echo "  make help       Show this help message"
	@echo ""
	@echo "First time setup:"
	@echo "  1. make install"
	@echo "  2. make serve"
	@echo ""

install:
	@echo "Installing Jekyll dependencies..."
	gem install -n /usr/local/bin bundler
	bundle install

serve:
	@echo "Starting Jekyll server..."
	@echo "Documentation will be available at:"
	@echo "  http://localhost:4000/jira-ticket-creator/"
	@echo ""
	bundle exec jekyll serve --baseurl /jira-ticket-creator

build:
	@echo "Building static site..."
	bundle exec jekyll build --baseurl /jira-ticket-creator

clean:
	@echo "Removing build artifacts..."
	rm -rf _site
	@echo "Done"
