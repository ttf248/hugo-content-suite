---
categories:
    - Computer
date: 2020-03-29T02:11:33+08:00
draft: false
featured: false
lastmod: 2020-03-29T02:11:33+08:00
slug: auto-integration-system-switch
subtitle: travis 到 github action
summary: ""
tags:
    - travis
    - github action
    - ci
    - blog
title: Automatic System Switchover
---


```bash
Test Code Block
```

- Apple
- Banana
- Orange
1. First step: Open document
2. Second step: Enter content
3. Third step: Save file

# Markdown Syntax Usage Full Breakdown

## One. Heading

Headings are used to structure documents. Markdown supports six levels of headings, achieved by using the '#' symbol followed by a space, with the level determined by the number of '#' symbols (the fewer '#' symbols, the higher the heading level).

### Grammar Examples:

### Display Effect:

# Heading (Highest Level)

## Level 2 Heading

### Level 3 Heading

#### Level 4 Heading

##### Level 5 Heading

###### Level 6 Heading (Lowest Level)

## II. Text Formatting

### 1. Bold (Bold)
**Purpose:** To highlight key information
**Syntax:** Enclose text within `**` tags
**Effect:** This is **bold text**, used to emphasize important points.

### 2. Italic (斜体)
**Purpose:** Indicates quotations, titles, or emphasis.
**Syntax:** Enclose text with * or _ characters.
**Effect:** *Italic text* is commonly used for quotations, _also available using underscores_ to achieve the same effect.

### 3. Bold and Italics (Bold + Italic)
**Purpose**: Applies both bold and italic effects.
**Syntax**: Enclose text with triple asterisks (`***`).
**Effect**: Text enclosed in `***simultaneously bold and italic***`.

### 4. Strikethrough
**Purpose:** To mark deleted or incorrect content
**Syntax:** Enclose text with `~~`
**Effect:** This is ~~strikethrough text~~, indicating that the content is invalid.

## III. Lists

### 1. Unordered List
**Purpose:** To display items with equal relationships
**Syntax:** Starts with `-`, `+`, or `*` followed by a space
- Apple
- Banana
- Orange

### 2. Ordered List
**Purpose:** To represent a sequence or steps
**Syntax:** Use numbers followed by a period (`.`) and a space (numbers can be repeated, Markdown will automatically render as continuous numbering)
1. First step: Open the document
2. Second step: Enter content
3. Third step: Save the file

### 3. Nested List
**Purpose:** Represents hierarchical relationships
**Syntax:** Indent subitems by two spaces or one tab
- Item 1
  - Subitem 1.1
  - Subitem 1.2
- Item 2
  1. Subitem 2.1
  2. Subitem 2.2

## Four. Quote (Blockquote)
**Purpose:** To cite others' viewpoints or literature
**Syntax:** Start with `>` followed by a space (can be nested multiple levels, adding one `>` for each level)
> This is a first-level quote.
> > This is a nested second-level quote.

## Five. Code

### 1. Inline Code
**Purpose:** To mark small amounts of code or terminology
**Syntax:** Enclose text with backticks (`` ` ``)
**Effect:** Output a greeting like `print("Hello World")`.

### 2. Code Block
**Purpose:** To display large blocks of code, supporting syntax highlighting.
**Syntax:** Enclose the code within backticks (`` ` ``), optionally specifying the programming language at the beginning (e.g., `python`, `javascript`).
**Effect (Python Syntax Highlighting):**

```python
def greet():
    print("Hello World")
greet()
```

## 6. Links

### 1. Hyperlink
**Purpose:** Navigate to a webpage or document
**Syntax:** `[Link Text](URL)`
**Effect:** Access [Baidu](https://www.baidu.com)

### 2. Reference Links
**Purpose:** Separates link content from text, suitable for long links.
**Syntax:**
**Effect:** View detailed explanation [Click Here][doc]

## Seven. Images

**Purpose:** Insert images

**Syntax:** `![Alt Text](Image URL "Optional Title")`

**Effect (if image exists):** ![Markdown Icon](https://example.com/markdown-icon.png "Markdown Logo")

## Eight. Tables (Tables)
**Purpose:** To present data in a structured manner.
**Syntax:**
- Use `|` to separate columns, and `-` to separate the header row from the content.
- A separator line (e.g., `---`, `:---`, or `:-:`) must be added below the header row.

| Zhang San | 25 | Engineer |

## VIII. Tables

| Name | Age | Occupation |
|---|---|---|
| Li Si | 30 | Designer |

## 9. Horizontal Rule
**Purpose:** To separate document sections
**Syntax:** Use `---`, `***`, or `___` (must be on a single line, with blank lines before and after)
**Example:**
**Result:**
Content ends

---
New section begins