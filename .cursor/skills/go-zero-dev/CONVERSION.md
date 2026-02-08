# Rules to Skills Conversion

This document tracks the conversion of `.cursor/rules/*.mdc` files to the `go-zero-dev` skill.

## Conversion Mapping

| Original Rule | Location in Skill | Status |
|---|---|---|
| `api.mdc` | SKILL.md - API Design section | ✅ Converted |
| `basics.mdc` | SKILL.md - Quick Reference section | ✅ Converted |
| `bool-fields.mdc` | SKILL.md - Type System Rules section | ✅ Converted |
| `cache-architecture.mdc` | SKILL.md - Cache Architecture section + examples.md | ✅ Converted |
| `code-reuse.mdc` | SKILL.md - Code Reuse section + examples.md | ✅ Converted |
| `db.mdc` | SKILL.md - Database Operations section + examples.md | ✅ Converted |
| `error-handling.mdc` | SKILL.md - Error Handling section + examples.md | ✅ Converted |
| `logging.mdc` | SKILL.md - Logging section + examples.md | ✅ Converted |
| `serialization.mdc` | SKILL.md - Type System Rules section | ✅ Converted |
| `style.mdc` | SKILL.md - Code Style section | ✅ Converted |

## File Structure

```
.cursor/skills/
├── README.md                          # Skill directory overview
└── go-zero-dev/
    ├── SKILL.md                      # Main skill (all rules consolidated)
    ├── examples.md                   # Detailed code examples
    └── CONVERSION.md                 # This file
```

## What Changed

### Organization

**Before (Rules):**
- 10 separate `.mdc` files
- Flat structure, each file independent
- No hierarchy or cross-references

**After (Skill):**
- Single skill with clear sections
- Progressive disclosure (SKILL.md → examples.md)
- Organized by topic with cross-references

### Discoverability

**Before (Rules):**
- Always loaded, no conditional application
- No metadata about when to use

**After (Skill):**
- Description tells agent when to apply
- Frontmatter metadata for better indexing
- Agent can choose to read or skip based on context

### Structure

**Before (Rules):**
```
.cursor/rules/
├── api.mdc
├── basics.mdc
├── bool-fields.mdc
└── ...
```

**After (Skill):**
```
.cursor/skills/go-zero-dev/
├── SKILL.md        # Core rules (~450 lines)
└── examples.md     # Detailed examples (~350 lines)
```

## Key Improvements

1. **Consolidated knowledge** - All Go-Zero best practices in one place
2. **Better structure** - Logical sections with clear hierarchy
3. **Progressive disclosure** - Quick reference in SKILL.md, details in examples.md
4. **Actionable checklists** - Added comprehensive checklist at the end
5. **Concrete examples** - Expanded examples section with complete code
6. **Cross-references** - Related topics linked together

## Migration Path

### Option 1: Keep Both (Recommended for transition)

- Keep existing rules active
- Test skill in parallel
- Gradually remove rules as skill proves effective

### Option 2: Remove Rules

```bash
# Backup first
mv .cursor/rules .cursor/rules.backup

# Test with just skills
# If issues arise, restore: mv .cursor/rules.backup .cursor/rules
```

### Option 3: Gradual Migration

1. Start with `basics.mdc` → test skill behavior
2. Remove `cache-architecture.mdc` → verify caching rules still work
3. Continue one-by-one until all rules removed

## Testing the Skill

Try these scenarios to verify the skill works:

1. **API Design**: Ask to create a new API endpoint
   - Should include @doc annotations
   - Should use proper field comments

2. **Database Operations**: Ask to add a new repository method
   - Should follow Repository pattern
   - Should include caching logic

3. **Code Reuse**: Ask to implement similar functionality
   - Should suggest extracting to helper
   - Should reference existing Logic

4. **Error Handling**: Ask to handle external API calls
   - Should distinguish 404 from other errors
   - Should use proper error checking

## Feedback

If you notice the skill missing something from the original rules:

1. Check if it's in examples.md
2. If not, edit SKILL.md to add it
3. Keep additions concise (remember: under 500 lines)
4. Use progressive disclosure for details

## Rollback

If you need to rollback to rules:

```bash
# The original rules are still in .cursor/rules/
# Just delete or rename the skill:
mv .cursor/skills/go-zero-dev .cursor/skills/go-zero-dev.disabled
```

The rules are still active, so removing the skill will revert to rule-based guidance.
