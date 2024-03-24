### Description de l'Application Web - Gestion de Deck Yu-Gi-Oh!

Cette application web tire parti de l'API Yu-Gi-Oh! pour fournir aux utilisateurs une plateforme de gestion de deck. L'API Yu-Gi-Oh! utilisée est disponible à l'adresse suivante : [API Yu-Gi-Oh!](https://ygoprodeck.com/api-guide/).

#### Fonctionnalités Principales

1. **Filtrage des Cartes**
   - Les utilisateurs peuvent filtrer les cartes en fonction de plusieurs critères, notamment la catégorie, le niveau et l'attribut. Ils peuvent spécifier ces critères dans l'URL de la page dédiée.

2. **Recherche de Cartes**
   - Les utilisateurs peuvent effectuer une recherche de cartes en saisissant un terme dans la barre de recherche. L'application renverra les cartes correspondant au terme de recherche dans leur nom.

3. **Gestion du Deck**
   - Une page "My Deck" est disponible pour permettre aux utilisateurs de gérer leur propre deck. Ils peuvent ajouter ou supprimer des cartes de leur deck, qui sera ensuite affiché dans cette page.

4. **Page d'Information sur une Carte**
   - Chaque carte a sa propre page d'information accessible en cliquant sur son image ou son titre. Cette page fournit des détails supplémentaires sur la carte sélectionnée.

5. **À Propos de Nous**
   - Une page "À propos de nous" est disponible pour fournir des informations sur l'application et son équipe de développement.

#### Endpoints de l'API Utilisés

Plusieurs endpoints de l'API Yu-Gi-Oh! sont utilisés pour récupérer les données nécessaires à l'application :

- **Endpoint 1:** https://db.ygoprodeck.com/api/v7/cardinfo.php
  - Récupère toutes les cartes disponibles dans la base de données.
- **Endpoint 2:** https://db.ygoprodeck.com/api/v7/cardinfo.php?id=
  - Récupère une carte spécifique en fonction de son ID.
- **Endpoint 3:** https://db.ygoprodeck.com/api/v7/cardinfo.php?fname=
  - Récupère une carte en fonction d'un terme de recherche dans son nom.

#### Pages Disponibles

Voici une liste des pages disponibles dans l'application :

- **Page d'accueil :** [http://localhost:8080/main](http://localhost:8080/main)
- **Liste des cartes :** [http://localhost:8080/liste](http://localhost:8080/liste)
- **Gestion du Deck :** [http://localhost:8080/deck](http://localhost:8080/deck)
- **Page de filtrage des cartes :** [http://localhost:8080/categorie](http://localhost:8080/categorie?categorie=Synchro+Monster&level=12)
- **Page d'information sur une carte :** [http://localhost:8080/info/"id"](http://localhost:8080/info/21123811)
- **Page de recherche de cartes :** [http://localhost:8080/recherche?query=""](http://localhost:8080/recherche?query=cosmic)
- **Page "À propos de nous" :** [http://localhost:8080/aboutus](http://localhost:8080/aboutus)

  
### Installation de l'Application

Pour installer l'application, suivez ces étapes :

1. Ouvrez votre IDE
2. Clonez : https://github.com/EnzoCarbo/carbonara_enzo_groupie_tracker
3. Lancez le main : go run .
4. Ouvrez vontre naviageur et utiliser le main présent dans la liste des pages disponibles.
5. Naviguez bien !

