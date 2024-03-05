# Projet de Groupe en Go :

# Système de gestion de plannings

<br>

<!-- toc -->

- [Objectif](#objectif)
- [Format du Projet](#format-du-projet)
- [Description détaillée et étapes du projet](#description-detaillee-et-etapes-du-projet)
  * [Interface Utilisateur en Console](#interface-utilisateur-en-console)
    + [Menu Principal](#menu-principal)
    + [Interaction Utilisateur](#interaction-utilisateur)
    + [Navigation](#navigation)
  * [Gestion du planning](#gestion-du-planning)
    + [Création d'Événements](#creation-devenements)
    + [Visualisation du planning](#visualisation-du-planning)
    + [Modification et Suppression](#modification-et-suppression)
  * [Base de Données pour les Événements](#base-de-donnees-pour-les-evenements)
    + [Structure de la Base de Données](#structure-de-la-base-de-donnees)
    + [Opérations sur la Base de Données](#operations-sur-la-base-de-donnees)
    + [Persistance et Configuration des Données](#persistance-et-configuration-des-donnees)
  * [Export du planning](#export-du-planning)
    + [Export JSON pour Sauvegarde](#export-json-pour-sauvegarde)
    + [Export CSV (Bonus)](#export-csv-bonus)
  * [Fonctionnalités Avancées](#fonctionnalites-avancees)
    + [Recherche et Filtres](#recherche-et-filtres)
    + [Rappels et Notifications](#rappels-et-notifications)
    + [Catégorisation du planning](#categorisation-du-planning)
    + [Interface utilisateur améliorée](#interface-utilisateur-amelioree)
  * [Bonus](#bonus)
    + [Fonctionnalités supplémentaires ou créatives au-delà des exigences de base.](#fonctionnalites-supplementaires-ou-creatives-au-dela-des-exigences-de-base)
  * [Suggestions de conception](#suggestions-de-conception)
    + [Bonnes Pratiques en Go](#bonnes-pratiques-en-go)
    + [Modularité et Structure du Code](#modularite-et-structure-du-code)
    + [Librairies Standard Utiles](#librairies-standard-utiles)
- [Récapitulatif](#recapitulatif)
  * [Récapitulatif des Étapes Clés](#recapitulatif-des-etapes-cles)
- [Rendu et Travail en Groupe](#rendu-et-travail-en-groupe)

<!-- tocstop -->

<div class="page-break"></div>

## Objectif

Créer un système complet de gestion de plannings en Go, intégrant une interface console interactive et une gestion
avancée des données.

## Format du Projet

- Groupes de 2 à 3 étudiants.
- Durée : 1 mois.

## Description détaillée et étapes du projet

### Interface Utilisateur en Console

#### Menu Principal

- Créez un menu principal clair et intuitif. Il doit offrir des options pour chaque fonctionnalité principale du système
  de gestion.
- **Exemple de menu** :

  ```markdown
  Bienvenue dans le Système de gestion de plannings
  -----------------------------------------------------
  1. Créer un nouvel événement
  2. Visualiser les événements
  3. Modifier un événement
  4. Supprimer un événement
  5. Rechercher un événement
  6. Quitter
  Choisissez une option : 
  ```

#### Interaction Utilisateur

- Gérez les entrées utilisateur avec soin. Assurez-vous que le système réagit de manière appropriée aux entrées valides
  et gère les erreurs ou saisies incorrectes sans crasher.
- **Exemple de gestion des entrées** :
    - Si l'utilisateur saisit un chiffre non valide, affichez un message d'erreur et demandez de nouveau une entrée.
    - Utilisez des boucles et des conditions pour gérer la navigation entre les menus.

#### Navigation

- Permettez à l'utilisateur de naviguer facilement entre les différentes fonctionnalités. Par exemple, après avoir
  visualisé ou créé un événement, l'utilisateur doit pouvoir retourner au menu principal ou effectuer une autre action
  sans redémarrer le programme.
- **Exemple de navigation** :
    - Après l'affichage du planning :

  ```markdown
  1. Retourner au menu principal
  2. Quitter
  Choisissez une option :
  ```

<div class="page-break"></div>

### Gestion du planning

#### Création d'Événements

- **Saisie des Détails** : Demandez à l'utilisateur de saisir tous les attributs nécessaires pour un nouvel événement
  (titre, date, heure, lieu, catégorie, description).
- **Validation des Données** : Assurez la validation des entrées (par exemple, vérifiez que la date et l'heure sont
  futures et formatées correctement).
- **Exemple de Saisie d'Événement** :

  ```text
  Entrez le titre de l'événement: Atelier Go
  Entrez la date (YYYY-MM-DD): 2023-12-05
  Entrez l'heure (HH:MM): 14:00
  Entrez le lieu: Salle A
  Choisissez une catégorie (Professionnel, Personnel, Loisir): Professionnel
  Entrez une brève description: Atelier d'initiation au langage Go.
  ```

#### Visualisation du planning

- **Affichage Global** : Implémentez une fonction pour afficher tous les événements, avec la possibilité de trier ou de
  filtrer la liste (par date, catégorie, etc.).
- **Affichage d'un Événement Spécifique** : Permettez aux utilisateurs de sélectionner et d'afficher les détails d'un
  seul événement, par exemple en entrant son identifiant ou son titre.
- **Exemple d'Affichage** :

  ```yaml
  Liste du planning :
  1. Atelier Go - 2023-12-05 - 14:00 - Professionnel
  2. Rencontre Sportive - 2023-12-12 - 10:00 - Loisir
  ...
  Entrez le numéro de l'événement pour voir plus de détails ou 0 pour revenir :
  ```

#### Modification et Suppression

- **Sélection d'Événements** : Permettez la sélection d'un événement spécifique à modifier ou supprimer, par exemple
  par son identifiant.
- **Processus de Modification** : Après la sélection d'un événement, offrez la possibilité de modifier chacun de ses
  attributs.
- **Confirmation de Suppression** : Avant de supprimer un événement, demandez une confirmation à l'utilisateur pour
  éviter les suppressions accidentelles.

<div class="page-break"></div>

### Base de Données pour les Événements

#### Structure de la Base de Données

- **Modélisation des Données** :
    - Utilisez une structure `Event` avec un identifiant entier. Ce dernier sera auto-incrémenté à chaque ajout d'un
      nouvel événement.
    - Stockez les événements dans une map avec l'ID comme clé pour un accès facile.
    ```go
    type Event struct {
        // ...
    }
    var eventsMap = make(map[int]Event)
    ```

#### Opérations sur la Base de Données

- **Ajout d'Événements** : Chaque nouvel événement se voit attribuer un ID unique qui s'auto-incrémente. Assurez-vous
  que
  l'ID est unique et n'a pas été utilisé auparavant.
- **Recherche d'Événements** : Implémentez des fonctions pour trouver un événement par son ID ou d'autres critères comme
  la date ou la catégorie.
- **Mise à Jour et Suppression** : Permettez la modification et la suppression d'événements en utilisant leur ID.

#### Persistance et Configuration des Données

- **Chargement des Identifiants** : Les identifiants de la base de données doivent être chargés depuis un fichier de
  configuration externe pour permettre des modifications dynamiques sans recompilation.
- **Fichier d'Initialisation de la Base de Données** : Créez un script ou un fichier d'initialisation pour préparer la
  base de données à son premier usage.

### Export du planning

#### Export JSON pour Sauvegarde

- Créez une fonction pour exporter les données actuelles du planning (la liste des événements) en format JSON. Cela peut être utilisé pour la
  sauvegarde ou l'export des données.
- Assurez-vous que l'export JSON reflète fidèlement la structure et l'état actuel de la base de données.

#### Export CSV (Bonus)

- En tant que fonctionnalité bonus, ajoutez une option pour exporter les données en format CSV, utile pour l'analyse de
  données ou la migration.

### Fonctionnalités Avancées

#### Recherche et Filtres

- **Implémentation de la Recherche** : Développez des fonctions pour rechercher du planning par titre, date, ou
  catégorie. Ces fonctions doivent parcourir la base de données et retourner les événements correspondant aux critères
  spécifiés.
- **Exemple d'Interface de Recherche** :

  ```yaml
  Entrez le titre pour la recherche: Atelier
  Résultats trouvés :
  - Atelier Go - 2023-12-05
  - Atelier Python - 2023-11-20
  ```

<div class="page-break"></div>

#### Rappels et Notifications

- **Système de Rappels** : Implémentez un mécanisme pour alerter les utilisateurs d'événements à venir, basé sur la
  date et l'heure courantes.
- **Affichage des Rappels** : Les rappels pourraient être des messages dans la console, déclenchés à l'ouverture du
  programme ou à des moments spécifiques.

#### Catégorisation du planning

- **Gestion des Catégories** : Permettez aux utilisateurs d'assigner des catégories aux événements (Professionnel,
  Personnel, Loisir, etc.) et de filtrer les événements en fonction de ces catégories.
- **Exemple d'une fonction de recherche par Catégorie** :

  ```go
  func searchEventsByCategory(eventsMap map[int]Event, category string) []Event {
      // Retournez les événements correspondant à la catégorie spécifiée
  }
  ```

#### Interface utilisateur améliorée

- **Interface Utilisateur Améliorée** : Améliorez l'interface utilisateur avec des fonctionnalités supplémentaires
  comme des menus déroulants, des couleurs, des commandes, etc.

### Bonus

> #### Fonctionnalités supplémentaires ou créatives au-delà des exigences de base.

### Suggestions de conception

#### Bonnes Pratiques en Go

- **Clarté et Simplicité** : Écrivez un code clair et lisible. Go privilégie la simplicité ; évitez les constructions
  compliquées inutiles.
- **Nommage des Variables et Fonctions** : Utilisez des noms descriptifs et cohérents pour les variables, fonctions et
  autres identifiants.
- **Commentaires et Documentation** : Documentez votre code avec des commentaires pertinents. Expliquez le 'pourquoi'
  derrière les choix de conception complexes.

#### Modularité et Structure du Code

- **Découpage en Packages** : Organisez votre code en packages logiques. Par exemple, un package pour la gestion de la
  base de données, un autre pour l'interface utilisateur, etc.
- **Fonctions Réutilisables** : Créez des fonctions qui accomplissent des tâches spécifiques et peuvent être réutilisées
  dans différents contextes du projet.

#### Librairies Standard Utiles

- **fmt pour l'Affichage et la Saisie** : Utilisez la librairie `fmt` pour afficher des informations à l'utilisateur et
  lire les saisies.
- **Gestion du Temps avec time** : La librairie `time` est essentielle pour manipuler les dates et les heures,
  notamment pour la création et le rappel des événements.
- **Manipulation de Fichiers avec os et io** : Pour la lecture et l'écriture des fichiers de configuration, de
  sauvegarde JSON, et d'export CSV.
- **Encodage JSON avec encoding/json** : Utilisez `encoding/json` pour convertir les données de la base de données en
  JSON et vice-versa.

<div class="page-break"></div>

## Récapitulatif

Vous avez désormais un aperçu complet du projet de création d'un système de gestion de plannings en Go. Ce projet
vous offre l'opportunité de mettre en pratique les compétences acquises en programmation Go, tout en relevant des défis
réalistes liés à la gestion de données, à l'interface utilisateur, et à la conception de logiciels.

### Récapitulatif des Étapes Clés

1. **Interface Utilisateur en Console** : Création d'une navigation intuitive et gestion des entrées utilisateur.
2. **Gestion du planning** : Ajout, visualisation, modification et suppression d'événements, avec une fonctionnalité
   pour afficher un événement spécifique.
3. **Base de Données et Gestion des Données** : Simulation d'une base de données en mémoire, persistance des
   identifiants et export de données en JSON et CSV.
4. **Fonctionnalités Avancées** : Recherche, rappels, catégorisation du planning, et export CSV bonus.
5. **Suggestions de Conception et Bonnes Pratiques** : Organisation du code, utilisation efficace des librairies
   standard de Go, et mise en œuvre de bonnes pratiques de programmation.

## Rendu et Travail en Groupe

- Code source complet et commenté.
- Fichier **README** avec instructions, description des fonctionnalités et répartition des tâches.
- Documentation technique décrivant la structure du code, les choix de conception, et la logique.

> Assurez-vous que le code est bien organisé, documenté, et testé, avant la soumission finale.

> Préparez une présentation détaillée du projet pour la soutenance, en mettant en avant les fonctionnalités
> implémentées, la collaboration au sein du groupe et les contributions individuelles.
