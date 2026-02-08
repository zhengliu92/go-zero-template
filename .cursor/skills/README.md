# Go-Zero Template Skills

This directory contains Cursor Agent Skills for the Go-Zero template project.

## Available Skills

### go-zero-dev

Comprehensive development guidelines for Go-Zero projects covering:

- **Architecture**: Layered architecture, Repository pattern, cache management
- **API Design**: API file structure, @doc annotations, field comments
- **Database Operations**: Repository pattern, helper functions, field updates
- **Caching**: Cache architecture rules, Repository-level caching
- **Logging**: Writer vs standard logging, log_type values, trace formatting
- **Code Reuse**: DRY principle, code extraction, calling existing Logic
- **Error Handling**: HTTP response errors, 404 handling
- **Type System**: Bool fields with `*bool`, JSON serialization with sonic
- **Code Style**: No emojis, no auto-documentation

## Usage

The Cursor agent will automatically apply the `go-zero-dev` skill when:

- Working with Go-Zero projects
- Designing APIs
- Implementing repository patterns
- Managing caches
- Handling database operations

## Files

- `go-zero-dev/SKILL.md` - Main skill instructions (quick reference and rules)
- `go-zero-dev/examples.md` - Detailed code examples and patterns

## Migration from Rules

This skill was converted from the `.cursor/rules/` directory. The original rules are still active but will eventually be deprecated in favor of this skill-based approach.

### Benefits of Skills over Rules

1. **Better organization** - Related content grouped together
2. **Progressive disclosure** - Main file is concise, detailed examples in separate file
3. **Better discoverability** - Description helps agent know when to apply
4. **Clearer structure** - Frontmatter metadata + markdown body

## Development

To modify this skill:

1. Edit `go-zero-dev/SKILL.md` for core rules and quick reference
2. Edit `go-zero-dev/examples.md` for detailed examples
3. Keep SKILL.md under 500 lines for optimal performance
4. Test changes by working on Go-Zero code and observing agent behavior

## Related Commands

Use Cursor commands (Ctrl+K or Cmd+K) for specialized tasks:

- `check-cache` - Verify cache architecture compliance
- `fix-cache` - Fix cache architecture violations
- `guard` - Add permission checks to endpoints
- `opt` - Check for code optimization opportunities
