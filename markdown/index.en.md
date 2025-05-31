---
categories:
    - Computer
date: 2020-03-29T02:11:33+08:00
draft: false
featured: false
lastmod: 2020-03-29T02:11:33+08:00
slug: auto-integration-system-switch
subtitle: Travis to GitHub Actions
summary: ""
tags:
    - travis
    - github action
    - ci
    - blog
title: Automatic System Switchover
---


- Apple
- Banana
- Orange
1. Step 1: Open the document
2. Step 2: Enter content
3. Step 3: Save the file

# Markdown Tag Usage Explained

## I. Heading

Headings are used to divide the document structure, and Markdown supports six levels of headings, implemented with a `#` symbol followed by a space. The level is determined by the number of `#` symbols (the fewer `#`, the higher the level).

### Syntax Example:

### Display Effect:

# Level 1 Heading (Top Level)

## Secondary Title

### Level 3 Heading

#### Level Four Heading

##### Level 5 Heading

###### Level 6 Heading (Smallest Level)

## II. Text Formatting

### 1. Bold

**Purpose:** To highlight key information.

**Syntax:** Wrap text with `**`.

**Effect:** This is **bold text**, used to emphasize important points.

### 2. Italic

**Purpose:** To indicate quotations, titles, or emphasize tone.  
**Syntax:** Wrap text with `*` or `_`.  
**Effect:** *Italicized text* is often used for citations, _also usable with underlines_.

### 3. Bold + Italic

**Purpose:** Overlays bold and italic effects.
**Syntax:** Wrap text in `***`.
**Effect:** Text that is ***both bold and italicized***.

### 4. Strikethrough
**Purpose**: Marks deleted or incorrect content
**Syntax**: Wrap text with ``~~``
**Effect**: This is ~~strikethrough text~~, indicating the content is invalid.

## III. Lists

### 1. Unordered List

**Purpose:** Displays entries with a parallel relationship.

**Syntax:** Starts with `-`, `+`, or `*` (recommend using `-`), followed by a space.

- Apple
- Banana
- Orange

### 2. Ordered List
**Purpose**: Represents sequence or steps
**Syntax**: Use a number followed by a `.` and a space (numbers can be repeated, Markdown will automatically render as consecutive numbers)
1. Step 1: Open the document
2. Step 2: Enter content
3. Step 3: Save the file

### 3. Nested List
**Purpose**: Represents hierarchical relationships
**Syntax**: Indent with 2 spaces or 1 tab before each item
- Main Item 1
  - Subitem 1.1
  - Subitem 1.2
- Main Item 2
  1. Subitem 2.1
  2. Subitem 2.2

## IV. Citations (Blockquote)
**Purpose:** To cite other people's views or literature.
**Syntax:** Starts with a `>` followed by a space (multiple levels of nesting are possible, each level adding another `>`).
> This is a first-level citation.
> > This is a nested second-level citation.

## V. Code

### 1. Inline Code

**Purpose**: Marks small snippets of code or terminology.

**Syntax**: Wrap text with `` ` ``.

**Effect**: Use `print("Hello World")` to output a greeting.

### 2. Code Block

**Purpose:** Displays large blocks of code and supports syntax highlighting.

**Syntax:** Wrap the code with ``` (you can specify a programming language at the beginning, e.g., `python`, `javascript`).

**Effect** (Python Syntax Highlighting):

## Six, Links

### 1. Hyperlink

**Function:** Jumps to a webpage or document

**Syntax:** `[Link Text](URL)`

**Example:** Visit [Baidu](https://www.baidu.com)

### 2. Reference Links

**Purpose:** Separates link content from text, suitable for long links.

**Syntax:**

**Effect:** View detailed explanation [Click here][doc]

## VII. Images

**Purpose:** Inserting images

**Syntax:** `![Alt text](Image URL "Optional Title")`

**Effect** (if the image exists): ![Markdown Icon](https://example.com/markdown-icon.png "Markdown Logo")

## VIII. Tables
**Purpose:** Structured data presentation
**Syntax:**
- Use `|` to separate columns, and `-` to separate the header from the content.
- A separator line must be added below the header row (`---` for left alignment, `:---` for right alignment, `:-:` for center alignment).
| Zhang San | 25 | Engineer |

## Eight, Tables

| Name | Age | Occupation |
| --- | --- | --- |
| Li Si | 30 | Designer |

## Nine, Horizontal Rule

**Purpose:** Separates document chapters.

**Syntax:** Use `---`, `***`, or `___` (must be on a separate line with blank lines before and after).

**Example:**

**Effect:**
Content ends

---
New chapter begins