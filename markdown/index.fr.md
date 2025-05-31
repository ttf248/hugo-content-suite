---
categories:
    - Ordinateur
date: 2020-03-29T02:11:33+08:00
draft: false
featured: false
lastmod: 2020-03-29T02:11:33+08:00
slug: auto-integration-system-switch
subtitle: travis vers GitHub Actions
summary: ""
tags:
    - travis
    - github action
    - ci
    - Blog
title: Commutation Intégrée Automatique
---


- Pomme
- Banane
- Orange
1. Étape 1 : Ouvrir le document
2. Étape 2 : Entrer du contenu
3. Étape 3 : Enregistrer le fichier

# Utilisation Complète des Marqueurs Markdown

## I. Titre (Heading)

Les titres servent à diviser la structure du document. Markdown prend en charge six niveaux de titre, mis en œuvre par le symbole `#` suivi d'un espace. Le niveau est déterminé par le nombre de `#` (`#` moins nombreux, niveau plus élevé).

### Exemple de Syntaxe :

### Effet Visuel :

# Titre Principal (Niveau Maximum)

## Titre Secondaire

### Titre de Niveau 3

#### Titre de Niveau 4

##### Titre de Niveau 5

###### Titre de niveau six (niveau minimal)

## II. Formatage du Texte (Text Formatting)

### 1. Gras (Bold)
**Fonctionnalité** : Mettre en évidence les informations clés.
**Syntaxe** : Entourez le texte avec `**`.
**Effet** : Ceci est du **texte en gras**, utilisé pour souligner les points importants.

### 2. Italique

**Fonctionnalité** : Indiquer une citation, un titre de livre ou mettre l'accent sur un ton.  
**Syntaxe** : Entourez le texte avec `*` ou `_`.
**Effet** : *Texte en italique* est couramment utilisé pour les citations, _l'italique peut également être réalisé avec des soulignements_.

### 3. Gras et Italique (Bold + Italic)
**Fonctionnalité** : Superposition des effets gras et italiques.
**Syntaxe** : Utiliser `***` pour englober le texte.
**Effet** : Texte *gras et italique* simultanément.

### 4. Barré (Strikethrough)
**Fonctionnalité** : Marquer du contenu supprimé ou erroné
**Syntaxe** : Entourez le texte avec `~~`
**Effet** : Ceci est ~~texte barré~~, indiquant que le contenu n'est plus valide.

## III. Listes

### 1. Liste Non Ordonnée (Unordered List)
**Fonctionnalité** : Afficher des éléments liés entre eux  
**Syntaxe** : Commencer par `-`, `+` ou `*` (recommandé `-`), suivi d'un espace
- Pomme
- Banane
- Orange

### 2. Liste Ordonnée (Ordered List)
**Fonctionnalité** : Représente une séquence ou des étapes
**Syntaxe** : Utilise un nombre suivi de `.` et d'un espace (les nombres peuvent être répétés, Markdown rendra automatiquement les numéros séquentiels)
1. Première étape : Ouvrir le document
2. Deuxième étape : Entrer du contenu
3. Troisième étape : Enregistrer le fichier

### 3. Liste imbriquée (Nested List)
**Fonctionnalité** : Représenter les relations hiérarchiques
**Syntaxe** : Indenter de 2 espaces ou d'une tabulation avant chaque élément enfant
- Élément principal 1
  - Élément enfant 1.1
  - Élément enfant 1.2
- Élément principal 2
  1. Élément enfant 2.1
  2. Élément enfant 2.2

## IV. Citations (Blockquote)
**Fonctionnalité** : Citer l'opinion ou la littérature d'autrui
**Syntaxe** : Commencer par un `>` suivi d’un espace (possibilité de créer des citations imbriquées en ajoutant un `>` pour chaque niveau)
> Ceci est une citation de premier niveau.
> > Ceci est une citation imbriquée de second niveau.

## V. Code

### 1. Code en Ligne (Inline Code)
**Fonctionnalité** : Marquer de petites quantités de code ou de termes techniques.
**Syntaxe** : Entourer le texte avec `` ` ``
**Exemple** : Utiliser `print("Hello World")` pour afficher un message d'accueil.

### 2. Bloc de Code (Code Block)
**Fonctionnalité** : Afficher des extraits de code importants, avec prise en charge de la coloration syntaxique.
**Syntaxe** : Entourez le code avec ``` (vous pouvez spécifier un langage de programmation au début, par exemple `python` ou `javascript`).
**Effet** (Coloration syntaxique Python) :

## VI. Liens

### 1. Lien Hypertexte (Hyperlink)
**Fonctionnalité** : Permet de se rendre vers une page web ou un document
**Syntaxe** : `[Texte du lien](URL)`
**Exemple** : Visitez [Baidu](https://www.baidu.com)

### 2. Liens de Référence

**Fonctionnalité** : Sépare le contenu du lien et le texte, idéal pour les longs liens.

**Syntaxe** :

**Résultat** : Pour plus d'informations, [cliquez ici][doc]

## VII. Images

**Fonctionnalité** : Insérer une image  
**Syntaxe** : `![Texte alternatif](URL de l'image "Titre optionnel")`  
**Effet** (si l'image existe) : ![Icône Markdown](https://example.com/markdown-icon.png "Logo Markdown")

## VIII. Tableaux (Tables)
**Fonctionnalité** : Présentation structurée des données
**Syntaxe** :
- Utiliser `|` pour séparer les colonnes et `-` pour séparer l'en-tête de table du contenu
- Ajouter une ligne de séparation sous la ligne d'en-tête (`---` pour alignement à gauche, `:---` pour alignement à droite, `:-:` pour centrage)
| Zhang San | 25 | Ingénieur |

## VIII. Tableaux (Tables)

| Nom | Âge | Profession |
|---|---|---|
| Li Si | 30 | Designer |

## IX. Ligne de Séparation (Horizontal Rule)
**Fonctionnalité** : Sépare les chapitres du document
**Syntaxe** : Utiliser `---`, `***` ou `___` (doit être sur une ligne séparée, avec des sauts de ligne avant et après)
**Exemple** :
**Résultat** :
Contenu terminé

---
Nouveau chapitre commencé