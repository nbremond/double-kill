# Double-kill

## Les spécifications

note:
Les modules Indexation et Recherche peuvent être regroupés selon l'architecture de la database.

## V1

### Indexation

+ à partir d'un répertoire donné, lister tous les sous fichiers et les ajouter à la database en spécifiant leur taille.

### Recherche

+ Retrouver dans la database les fichiers ayant le même taille, les comparer bit à bit et en cas d'égalité, l'indiquer dans la data base.


### Supression


## V2

### Indexation

+ mettre à jour dans la database les fichiers qui ont changé.
+ supprimer de la database les fichiers retirés du dossier.

### Recherche


### Supression

+ pour les fichiers en plusieurs exemplaires, demander à l'utilisateur le ou lequels suprimer.
+ Avant chaque suppression on refait une comparaison bit à bit
+ La base de donnée est mise à jour au fur et à meusure des supressions.


## V3

### Indexation

### Recherche

### Supression

+ l'utilisateur peut spécifier des dossiers préférentiels pour conserver ou supprimer les doublons.
+ ajout d'une option pour ne faire que les supressions automatiques (pour execution en arrière plan)
+ ajout d'une option pour ignorer la vérification bit à bit. (on verifie quand même l'existance des fichier avant de supprimer)


## V4

### Indexation

+ ajout d'un hash du fihier complet (contenu et en-tête) dans la base de donnée.
+ l'utilisateur spécifie des regles pour savoir quand le hash est calculé (selon la taille, selon le type,)

### Recherche

+ la recherche vérifie l'égalité des hash avant de faire l'égalité bit à bit
+ si une vérification bit à bit est nécéssaire, on calcule le hash en même temps *il faut définir si on poursuit le calcul du hash si la comparaison échoue au deuxième bit*

### Supression


## V5

### Indexation

+ Pour les types de fichier le permettant, ajouter la taille et un hash du contenu sans en-tête. (exemple : pour les mp3 on hash la musique mais pas les informations auteur, album, ...)

### Recherche

+ Il faut maintenant rechercher les fichiers différents mais qui ont le même contenu et faire la comparaison bit à bit du contenu
+ Lorsqu'une correspondance est trouvée, il faut la stocker en indiquant que l'en-tête est différent.
+ Bien évidement, il ne faut pas que les fichiers completement identiques soient marqués comme ayant le même contenu, ils doivent être marqués comme correspondance totale.

### Supression

+ la suppression des fichiers à contenu indentique se fait au cas par cas.


## V6

### Indexation

### REcherche

### Supression

+ créer des liens physiques au lieu de suprimer.
+ créer des liens symboliques au lieu de supprimer.


## V8

### Indexation
+ répertorier les dossiers vides(éventuelement déjà implémenté selon l'archi database)

### Recherche

### Supression

+ suprimer les dossiers vides.

## Ce qui doit être implémenté a terme. (soit implémentation au fur et à meusure soit implementation apres tout le reste)

mise en place d'un fichier de conf

### Indexation

+ On spécifie les dossiers à indexer en précisant si besoin les types de fichiers à ingorer, la taille min des fichiers à indexer, ...

### Recherche

### Supression

+ Chaque bloc d'option concerne un ou plusieurs dossiers.
+ On indique les dossiers à conserver ou a vider par ordre de préférence.
+ On indique les types des fichiers à ne pas supprimer ou au contraire une iste de ceux qu'on  veut supprimer.
    Example:

        global{
        ignore=/usr/bin,/etc;
        }

        type Musique{
        extension=mp3,wav;
        OR
        [extension=.jpg
        AND
        filename contain folder_cover]
        }

        delete Music{//est-ce utile de rajouter un nom aux règles ?
        type =Musique; 
        folders= ~Nicolas/, ~Fabrice/, /share/Musique;
        ignored_folders=.PlayonLinux;
        keep= /share/Musique;
        }

        make_hard_link RessourcesJeux{
        type.extension=jpg,jpeg,mp3,png;
        folders= ~Nicolas/.PlayOnLinux, ~Fabrice/.PlayOnLinux;
        ignored_folders=*conf*
        }





## V9
On peut utiliser plusieurs database afin de comparer son ordinateur et un serveur de fichier.
