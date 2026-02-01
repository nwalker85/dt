# dt - DEVONthink CLI

A command-line interface for DEVONthink automation on macOS. Wraps DEVONthink's AppleScript/JXA capabilities for scripting and workflow integration.

## Installation

```bash
# Build from source
git clone https://github.com/nwalker85/dt.git
cd dt
go build -o dt .

# Install (choose one)
sudo cp dt /usr/local/bin/
# or
ln -s $(pwd)/dt ~/.local/bin/dt
```

## Usage

```bash
dt <command> [flags]
```

Global flags:
- `--json` - Output in JSON format (where applicable)

## Core Commands

### Database Operations

```bash
dt stats                    # Item count per database
dt summary                  # Database details as JSON
dt databases                # List database names
```

### Search & Discovery

```bash
dt search <query>           # Search and return file paths
dt recent [--days N]        # Recently added/modified items
dt info <uuid>              # Detailed item information
dt duplicates               # Find duplicate items
```

### Tag Management

```bash
dt tags list                # List all unique tags
dt tag <query> <tags...>    # Add tags to matching items
dt untag <query> <tags...>  # Remove tags from matching items
```

### File Operations

```bash
dt open <query>             # Open matching items in DEVONthink
dt export <tag> [--dest]    # Export tagged items to filesystem
dt import <source> [--db]   # Import folder into database
dt move <query> --to <db>   # Move items to another database
dt trash <query>            # Move items to trash
```

### Content Operations

```bash
dt create <name> [--db]     # Create new document
dt ocr <query>              # OCR matching documents
dt classify <query>         # Auto-classify items using AI
dt see-also <uuid>          # Find related documents
```

## Workflow Commands

Higher-level operations that combine multiple actions:

```bash
dt inbox                    # List inbox items across all databases
dt inbox process            # Interactive inbox processing
dt archive <query>          # Tag as archived + move to Archive group
dt weekly-report            # Summary of items added this week
```

## DEVONthink Query Syntax

dt uses DEVONthink's native search syntax:

| Operator | Description | Example |
|----------|-------------|---------|
| `kind:` | File type | `kind:pdf`, `kind:markdown` |
| `tag:` | Tag | `tag:work`, `tag:important` |
| `name:` | Filename | `name:report` |
| `content:` | Content search | `content:quarterly` |
| `date:` | Date | `date:today`, `date:thisweek` |
| `size:` | File size | `size:>1mb` |

Boolean operators: `AND`, `OR`, `NOT`

### Examples

```bash
# Find all PDFs tagged "work"
dt search "kind:pdf tag:work"

# Export important documents
dt export important --dest ~/exports

# Tag all markdown files from this week
dt tag "kind:markdown date:thisweek" review

# Find related documents
dt see-also 12345-67890-ABCDE

# Process inbox
dt inbox process
```

## Scripting Examples

### Backup tagged items

```bash
#!/bin/bash
dt export backup --dest "/Volumes/Backup/devonthink/$(date +%Y-%m-%d)"
dt untag "tag:backup" backup
```

### Daily inbox report

```bash
#!/bin/bash
echo "Inbox items: $(dt inbox --json | jq length)"
dt recent --days 1 --json | jq -r '.[].name'
```

### Archive old items

```bash
#!/bin/bash
dt archive "date:<lastmonth tag:processed"
```

## Requirements

- macOS
- DEVONthink 3 (or DEVONthink)
- Go 1.21+ (for building)

## License

MIT
