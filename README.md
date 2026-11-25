# Push Swap en Go

Deux programmes travaillent sur les mêmes nombres. Avec `push-swap "3 1 2"`, tu demandes « quels mouvements permettent de trier ces trois nombres ? ». Avec `checker "3 1 2"`, tu demandes « est-ce que les mouvements que je fournis fonctionnent ? ».

## 1. Lancer le projet

Il te faut Go 1.18 ou plus récent. Décompresse l’archive, ouvre le dossier qui contient `go.mod` dans ton éditeur, puis ouvre un terminal dans ce dossier. Sur Mac ou Linux, tape les commandes suivantes, une ligne à la fois.

```sh
go version
go build -o push-swap ./cmd/push-swap
go build -o checker ./cmd/checker
go test ./...
```

Les deux commandes `go build` créent les exécutables à côté de `go.mod`. Tous les imports viennent du projet ou de la bibliothèque standard de Go.

Pour voir les mouvements, lance le premier programme seul.

```sh
./push-swap "3 1 2"
```

Il affiche `ra`. Le sommet de A est à gauche. En envoyant le 3 à la fin, on passe de `[3 1 2]` à `[1 2 3]` en un mouvement.

Pour vérifier automatiquement la réponse, utilise le caractère `|`. Le terminal transmet chaque ligne affichée par le premier programme au deuxième.

```sh
./push-swap "3 1 2" | ./checker "3 1 2"
```

Résultat attendu : `OK`. Les nombres sont écrits deux fois parce que les deux programmes démarrent chacun avec leur propre pile A.

Pour vérifier les deux exécutables et leurs sorties depuis le terminal, lance `sh verifier.sh` sur Mac ou Linux. Ce script compile dans un dossier temporaire, lance les tests Go, compare les sorties et rejoue les exemples du sujet. Une erreur interrompt le script.

## 2. Où lire le code ?

| Fichier | Ce que tu y trouveras |
| --- | --- |
| `go.mod` | Le nom du module, `pushswap`, utilisé pour importer `pushswap/piles`. |
| `piles/arguments.go` | La lecture des nombres et la détection des erreurs de saisie. |
| `piles/operations.go` | La représentation des piles, les quatre mouvements de base et les onze instructions. |
| `piles/chemin.go` | La recherche du chemin minimal pour les petites piles. |
| `piles/tri.go` | Le choix des mouvements et l’enregistrement des instructions dans l’ordre. |
| `cmd/push-swap/main.go` | Le point de départ du programme qui affiche la solution. |
| `cmd/checker/main.go` | Le point de départ du programme qui lit et vérifie les mouvements. |
| `piles/piles_test.go` et `cmd/checker/main_test.go` | Les cas de test. Tu peux ajouter tes propres exemples dans les listes existantes. |

Chaque dossier dans `cmd` contient sa propre fonction `main`, le point de départ d’un exécutable Go. Les deux importent le même paquet `piles`. Modifier le comportement de `sa` dans ce paquet modifie donc la façon dont les deux programmes l’exécutent.

## 3. Passer du texte à la pile A

Regarde `LireNombres` dans `arguments.go`. Quand tu écris `./push-swap "4 -2 7"`, le terminal transmet le texte `4 -2 7` en un seul argument. Avec `./push-swap 4 -2 7`, il transmet trois arguments. Les deux écritures donnent ici la même pile.

1. `os.Args[1:]` récupère les arguments après le nom du programme.
2. `strings.Fields(argument)` découpe chaque argument en mots en tenant compte des espaces.
3. `strconv.Atoi(mot)` convertit un mot en entier Go. Par exemple, `"-2"` devient le nombre `-2`. Une lettre, une virgule ou un nombre hors des limites de `int` provoque une erreur.
4. `dejaVus` est une table qui associe un entier à un booléen. Si `dejaVus[nombre]` est vrai, on a déjà rencontré ce nombre. Comme on vérifie après la conversion, `2` et `+02` comptent comme un doublon.
5. `append(nombres, nombre)` ajoute chaque valeur à la fin. Le premier nombre saisi reste au sommet.

Sans argument, les deux programmes se terminent silencieusement. Fournir explicitement une chaîne vide, comme `./push-swap ""`, déclenche `Error`.

## 4. Comprendre les piles et les instructions

Dans `operations.go`, `A` et `B` sont de type `[]int`. En Go, c’est une *slice*, une suite d’entiers dont la longueur peut changer. Si A vaut `[8 3 5]`, `A[0]` vaut 8 et `A[1:]` représente `[3 5]`. B commence vide.

| Instruction | Mouvement |
| --- | --- |
| `sa`, `sb` | Échanger les deux premiers éléments de A ou de B. |
| `ss` | Faire les deux échanges. |
| `pa` | Prendre le sommet de B et le poser au sommet de A. |
| `pb` | Prendre le sommet de A et le poser au sommet de B. |
| `ra`, `rb` | Déplacer le premier élément de A ou B vers la fin. |
| `rr` | Faire cette rotation sur les deux piles. |
| `rra`, `rrb` | Déplacer le dernier élément de A ou B vers le début. |
| `rrr` | Faire cette rotation inverse sur les deux piles. |

Le nom de `pa` ou `pb` indique la pile **qui reçoit** le nombre. Voici par exemple ce que fait `pousser` avec A = `[8 3 5]` et B = `[1 4]`. Pour `pb`, la source est A et la destination est B. La fonction garde le 8, renvoie `[3 5]` comme nouvelle source et `[8 1 4]` comme nouvelle destination. Dans `append([]int{sommet}, destination...)`, les trois points demandent d’ajouter les éléments de la destination un par un après le sommet.

`Appliquer` lit le nom du mouvement avec un `switch` et choisit la bonne fonction. Le récepteur `(p *Piles)` lui donne accès aux piles de l’appelant. Une instruction autorisée renvoie `true`, même si la pile est trop petite pour bouger. Un nom inconnu renvoie `false`. Avec `ss`, si A contient un seul nombre et B en contient deux, l’échange se fait quand même sur B.

## 5. Comment push-swap choisit ses mouvements

Commence par `Trier` dans `tri.go`. Si A est déjà triée, la liste d’instructions est vide. Sinon, la taille de A détermine la méthode.

### Jusqu’à six nombres

Dans `piles/chemin.go`, la fonction `cheminLePlusCourt` essaie tous les premiers mouvements, puis les suites de deux mouvements, puis celles de trois, et ainsi de suite. C’est une **recherche en largeur**. Chaque essai contient les deux piles et le chemin suivi. On copie les piles avant chaque mouvement pour que les essais restent indépendants. La table `dejaVus` évite de réexplorer une disposition déjà atteinte. Comme chaque mouvement coûte un, la première solution trouvée a forcément la longueur minimale.

Pour l’expliquer à l’oral, prends `[4 1 3 2]`. Le chemin retenu est `ra pb sa pa`. Le programme explore aussi les autres choix possibles. Aucun chemin en trois mouvements ne suffit. On limite cette recherche à six nombres parce que le nombre de dispositions augmente très vite.

| Exemple du sujet | Une solution minimale dans cette version | Nombre de mouvements |
| --- | --- | --- |
| `4 67 3 87 23` | `pb ra pb rr pa pa` | 6 |
| `2 1 3 6 5 8` | `pb pb ra sa rrr pa pa` | 7 |

### À partir de sept nombres

Le programme regarde d’abord si échanger les deux premiers suffit. Il regarde aussi si la pile est déjà triée en cercle. Par exemple, pour `[7 1 2 3 4 5 6]`, un seul `ra` suffit. Si la pile est triée en cercle, on ramène le minimum au sommet en choisissant le sens de rotation le plus court.

Pour les autres entrées, le programme utilise le tri suivant. La fonction locale `jouer` applique chaque mouvement et ajoute son nom à la liste des instructions à afficher.

`trierParBits` utilise un **tri radix binaire**. Un bit est un chiffre qui vaut 0 ou 1. On classe les nombres d’après leurs bits successifs, en commençant à droite.

Prenons A = `[2 -3 1 -2 0 -1 3]`. Le minimum vaut −3. Pour lire les bits, on utilise l’écart au minimum. Par exemple, pour 2, l’écart vaut 5, soit `101` en binaire. Pour −2, il vaut 1, soit `001`. Trier les écarts revient à trier les nombres, car on soustrait le même minimum à chacun. Les piles gardent les valeurs d’origine.

| Nombre | −3 | −2 | −1 | 0 | 1 | 2 | 3 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Écart | 0 | 1 | 2 | 3 | 4 | 5 | 6 |
| Écart en binaire | 000 | 001 | 010 | 011 | 100 | 101 | 110 |

À chaque passage, on traite autant de sommets qu’il y avait de nombres au début du passage. Pour un bit à 0, `pb` envoie le nombre dans B. Pour un bit à 1, `ra` envoie le nombre à la fin de A. Ensuite, on vide B avec `pa`. Les zéros reviennent devant les uns. Les deux transferts successifs inversent deux fois leur ordre, donc les zéros gardent entre eux le classement du passage précédent. Les uns le gardent aussi grâce aux rotations.

| Après quel passage ? | Contenu de A, du sommet vers le bas |
| --- | --- |
| Départ | 2, −3, 1, −2, 0, −1, 3 |
| Bit de droite, masque 1 | −3, 1, −1, 3, 2, −2, 0 |
| Bit du milieu, masque 2 | −3, 1, 2, −2, −1, 3, 0 |
| Bit de gauche, masque 4 | −3, −2, −1, 0, 1, 2, 3 |

**Lire les symboles dans le code.** `ecart & masque` regarde si le bit choisi est présent. `masque <<= 1` passe au bit suivant, donc les masques valent 1, 2, 4, 8… `bitsUtiles |= ecart` rassemble les positions où au moins un écart a un 1. Cela permet de sauter une position entièrement à zéro. La conversion en `uint` avant la soustraction permet de représenter tout l’écart entre deux entiers, même si le minimum est très négatif. La condition `masque != 0` arrête la boucle si le masque a dépassé le dernier bit disponible.

Cette méthode privilégie un tri explicable et s’arrête dès que A est triée. Elle peut produire davantage de mouvements que la séquence la plus courte possible.

## 6. Comment checker décide entre OK, KO et Error

`verifier` lit une ligne à la fois sur l’entrée standard. Le retour à la ligne termine l’instruction. Le programme retire ce séparateur, puis appelle la même méthode `Appliquer` que le trieur. Il garde les espaces éventuels pour détecter une instruction mal écrite.

| Situation | Résultat |
| --- | --- |
| A est triée et B est vide après la dernière instruction | `OK` |
| Les instructions sont valides, mais A est mal triée ou B contient encore un nombre | `KO` |
| Nombre invalide, doublon, instruction inconnue, ligne vide ou instruction mal terminée | `Error` sur la sortie d’erreur |
| Aucun argument | Aucune sortie |

Essaie ces commandes depuis le dossier du projet.

| Commande | Résultat attendu |
| --- | --- |
| `./push-swap "1 2 3"` | Aucune sortie. Les nombres sont déjà triés. |
| `printf 'sa\n' \| ./checker "2 1 3"` | `OK`. |
| `printf 'pb\n' \| ./checker "1 2 3"` | `KO`. A est triée, mais B contient encore 1. |
| `printf 'rotation\n' \| ./checker "2 1"` | `Error`. Ce mouvement n’existe pas. |
| `./push-swap "2 +02"` | `Error`. Les deux écritures représentent le même nombre. |

## 7. Vérifier que tu sais l’expliquer

Prends quatre papiers marqués 4, 1, 3 et 2. Pose-les dans cet ordre pour former A, avec le 4 au sommet. Joue à la main `ra`, `pb`, `sa`, puis `pa`. Après chaque mouvement, note les deux piles. À quel moment B reçoit-elle le 1 ? Pourquoi faut-il le remettre dans A même si les autres nombres sont déjà dans l’ordre ?

Tu peux vérifier ta réponse avec le programme.

```sh
printf 'ra\npb\nsa\npa\n' | ./checker "4 1 3 2"
```

Tu dois obtenir `OK`. Après `ra`, A vaut `[1 3 2 4]`. Après `pb`, B contient 1 et A vaut `[3 2 4]`. L’échange donne `[2 3 4]`, puis `pa` remet le 1 au sommet.
