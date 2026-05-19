# Bob's Coding Rules and Standards

## Markdown Formatting Rules

### MD022/blanks-around-headings

When generating or editing markdown files, always observe MD022/blanks-around-headings:

- Add a blank line before each heading (except at the start of the file)
- Add a blank line after each heading
- This improves readability and follows markdown best practices

**Example (Correct):**

```markdown
# Main Title

This is content after the heading.

## Section Heading

More content here.
```

**Example (Incorrect):**

```markdown
# Main Title
This is content without blank line.
## Section Heading
More content here.
```

### MD047/single-trailing-newline

When generating or editing markdown files, always observe MD047/single-trailing-newline:

- Files should end with a single newline character
- No multiple trailing newlines
- No missing trailing newline
- This ensures consistent file endings across different editors and systems

**Example (Correct):**

```markdown
# Document Title

Content here.
[single newline at end]
```

**Example (Incorrect):**

```markdown
# Document Title

Content here.[no newline]
```

or

```markdown
# Document Title

Content here.
[multiple newlines]



```

### MD032/blanks-around-lists

When generating or editing markdown files, always observe MD032/blanks-around-lists:

- Add a blank line before each list (ordered or unordered)
- Add a blank line after each list
- This improves readability and ensures proper rendering

**Example (Correct):**

```markdown
Here is some text.

- List item 1
- List item 2
- List item 3

More text after the list.
```

**Example (Incorrect):**

```markdown
Here is some text.
- List item 1
- List item 2
- List item 3
More text after the list.
```

### MD031/blanks-around-fences

When generating or editing markdown files, always observe MD031/blanks-around-fences:

- Add a blank line before each fenced code block
- Add a blank line after each fenced code block
- This improves readability and ensures proper rendering
- Applies to both ``` and ~~~ style fences

**Example (Correct):**

```markdown
Here is some text.

```bash
echo "Hello World"
```

More text after the code block.
```

**Example (Incorrect):**

```markdown
Here is some text.
```bash
echo "Hello World"
```
More text after the code block.
```
```


### MD036/no-emphasis-as-heading

When generating or editing markdown files, always observe MD036/no-emphasis-as-heading:

- Do not use emphasis (bold or italic) as a substitute for proper headings
- Use proper heading syntax (# ## ### etc.) instead of **Bold Text** or *Italic Text* for section titles
- This ensures proper document structure and accessibility
- Screen readers and document parsers rely on proper heading hierarchy

**Example (Correct):**

```markdown
## Section Title

This is the content of the section.

### Subsection Title

More content here.
```

**Example (Incorrect):**

```markdown
**Section Title**

This is the content of the section.

*Subsection Title*

More content here.
```


### MD004/ul-style

When generating or editing markdown files, always observe MD004/ul-style:

- Use a consistent unordered list marker style throughout the document
- Choose one style and stick with it: asterisk (*), plus (+), or dash (-)
- The most common and recommended style is dash (-)
- Mixing styles within the same document reduces readability and consistency

**Example (Correct - using dash):**

```markdown
- List item 1
- List item 2
  - Nested item
  - Another nested item
- List item 3
```

**Example (Correct - using asterisk):**

```markdown
* List item 1
* List item 2
  * Nested item
  * Another nested item
* List item 3
```

**Example (Incorrect - mixed styles):**

```markdown
- List item 1
* List item 2
  + Nested item
  - Another nested item
* List item 3
```

**Note:** While all three markers are valid markdown, consistency improves readability and maintainability.

### MD010/no-hard-tabs

When generating or editing markdown files, always observe MD010/no-hard-tabs:

- Use spaces instead of hard tabs for indentation
- Hard tabs can render inconsistently across different editors and viewers
- Most markdown parsers expect spaces for indentation
- Use 2 or 4 spaces for nested list items (be consistent within a document)

**Example (Correct):**

```markdown
- List item 1
  - Nested item (2 spaces)
  - Another nested item
- List item 2
    - Nested item (4 spaces)
    - Another nested item
```

**Example (Incorrect):**

```markdown
- List item 1
	- Nested item (hard tab)
	- Another nested item
- List item 2
		- Nested item (hard tabs)
		- Another nested item
```

**Note:** Configure your editor to convert tabs to spaces when editing markdown files.

### MD060/table-column-style

When generating or editing markdown files, always observe MD060/table-column-style:

- Use consistent column alignment across all rows in a table
- Align the separator row (dashes) with the header and data rows
- Pad cells with spaces to maintain visual alignment
- This improves readability and makes tables easier to maintain

**Example (Correct):**

```markdown
| Command                              | Description                      |
|--------------------------------------|----------------------------------|
| `mvn clean`                          | Remove all build artifacts       |
| `mvn compile`                        | Compile source code              |
| `mvn test`                           | Run unit tests                   |
| `mvn package`                        | Create JAR files                 |
| `mvn install`                        | Install to local Maven repository|
```

**Example (Incorrect):**

```markdown
| Command | Description |
|---------|-------------|
| `mvn clean` | Remove all build artifacts |
| `mvn compile` | Compile source code |
| `mvn test` | Run unit tests |
| `mvn package` | Create JAR files |
| `mvn install` | Install to local Maven repository|
```

**Note:** While the incorrect example will render correctly, the aligned version is more maintainable and easier to read in source form.
**Note:** Emphasis is fine for highlighting words within paragraphs, just not as standalone headings.
