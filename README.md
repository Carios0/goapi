# ToDoApp-API


## 12. Testing

For running a crude test, please execute the commands found in
**curlTest.sh**
in the root directory before inserting any items yourself (tests will fail if auto-generated Ids don't start at 1). In case you need to restart do ```docker-compose down``` and ```docker-compose up```. This performs some creation, updating and deletion tasks (for details see curlTest.sh). You should be left with 2 items. They differ in all field and only the item with id=3 has comments and any upvotes, so it is possible to manually test all GET requests and filterings offered with just these two.

## 11. Authentication

All functionalities except the home page require authentication by the user. The valid username is "leo" and the password is "123". 

## 1. Inserting Items

For inserting an item send a POST request to the endpoint /items like this:

```curl -X POST -H "Content-Type: application/json" -d '{"title": "Sample Item 1", "description": "First test item for the ToDo list.", "author": "Leo G."}' http://localhost:8080/items -u leo:123```

It is required to set header to "Content-Type: application/json" and to provide the arguments title, description and author in the json payload of which only description may be empty. All other field of the newly created item are set automatically (id, time) or by default (done, priority).

## 2. Fetching one Item

To fetch a single item send a GET request to the /items/{id} endpoint like this:

```curl localhost:8080/items/1 -u leo:123```

This will fetch the item with id=1.

## 3. Fetch all Items

To fetch all items send a GET request to the /items/all endpoint like this:

```curl localhost:8080/items/all -u leo:123```

## 4.&6. Update Item

To update an item send a PATCH request to the /items/{id} endpoint like this:

```curl -X PATCH "localhost:8080/items/3?field=title&value=Updated%20Title" -u leo:123```

This will return an error, if the field parameter is not a valid item field. Valid item fields are:

1. title
2. description
3. done (default: false)
4. priority (default: 2, possible values: 1 , 2 , 3 , 4)

All other fields cannot be updated by design, because we don't want someone manipulating ids, authors or timestamps.
Whitespaces in the value parameter are encoded with "%20".
Use this method to update the field "done" to mark an item as done ("true") or not done ("false").

## 5. Delete Item

To delete an item send a DELETE request to the endpoint /items/{id} like this:

```curl -X DELETE localhost:8080/items/1 -u leo:123```

## 7. Filter by Author

To get all items by any one author send a GET request to the endpoint /items/all/author/{author} like this:

```curl localhost:8080/items/all/author/Leo%20G -u leo:123```

## 8. Filter by Status (done)

To filter all items by their status, send GET request to /items/all/done/{done} like this:

```curl localhost:8080/items/all/done/true -u leo:123```

## 9. Sort by Time

To fetch all items sorted by time, send GET request to /items/all/sort/time/{inv} where {inv} lets you invert the result. Example:

```curl localhost:8080/items/all/sort/time/true -u leo:123```

## 10 Sort by Priority

To fetch all items sorted by priority (has to be manually set), send GET request to /items/all/sort/priority/{inv}. Again, {inv} inverts the sorting if set as "true". Example:

```curl localhost:8080/items/all/sort/priority/false -u leo:123```

## 13 Voting

To vote for an item, send POST request to /votes/ like this:

```curl -X POST -H "Content-Type: application/json" -d '{"itemId":3, "description": "Comment Number One", "user": "Carl G"}' http://localhost:8080/comments -u leo:123```

## 14 Rank Items by Vote

To see a list of items in order of vote count, send a GET request to /items/all/sort/votes/{inv} (inv = inverted) like this:

```curl localhost:8080/items/all/sort/votes/true```

## 15 Comments

To comment on an item use POST on /comments like this:

```curl -X POST -H "Content-Type: application/json" -d '{"itemId":3, "description": "Comment Number One", "user": "Carl G"}' http://localhost:8080/comments -u leo:123```

To delete a comment use DELETE on /comments/{id}, where id is the id of the comment. Example:

```curl -X DELETE localhost:8080/comments/3 -u leo:123```

To see all comments on one item, use GET on /comments/{itemId} like this:

```curl localhost:8080/comments/3```
