
#!/bin/bash
# insert item 1
curl -X POST -H "Content-Type: application/json" -d '{"title": "Sample Item 1", "description": "First test item for the ToDo list.", "author": "Carl G"}' http://localhost:8080/items -u leo:123
# insert item 2
curl -X POST -H "Content-Type: application/json" -d '{"title": "Delete this", "description": "This item should be deleted in the course of this test.", "author": "Leo G"}' http://localhost:8080/items -u leo:123
# wait for 2 seconds to ensure different timestamp on item 3
sleep 2
# insert item 3
curl -X POST -H "Content-Type: application/json" -d '{"title": "Update this", "description": "This item should be updated in the course of this test.", "author": "Leo G"}' http://localhost:8080/items -u leo:123
# update field description to "The update worked" on item 3
curl -X PATCH "localhost:8080/items/3?field=description&value=The%20Update%20worked.%20This%20item%20should%20have%20upvotes%20and%20comments." -u leo:123
# update field title to "Updated Title" on item 3
curl -X PATCH "localhost:8080/items/3?field=title&value=Updated%20Title" -u leo:123
# update field priority to 4 on item 3
curl -X PATCH "localhost:8080/items/3?field=priority&value=4" -u leo:123
# update field done to true on item 3
curl -X PATCH "localhost:8080/items/3?field=done&value=true" -u leo:123
# delete item number 2
curl -X DELETE localhost:8080/items/2 -u leo:123
# upvote item 3
curl -X POST -H "Content-Type: application/json" -d '{"itemId":3, "user": "Leo G"}' http://localhost:8080/votes -u leo:123
# upvote item 1
curl -X POST -H "Content-Type: application/json" -d '{"itemId":1, "user": "Carl G"}' http://localhost:8080/votes -u leo:123
# delete upvote on item 1
curl -X DELETE "localhost:8080/votes/1?user=Carl%20G" -u leo:123
# add 2 comments on item 3
curl -X POST -H "Content-Type: application/json" -d '{"itemId":3, "description": "Comment Number One", "user": "Carl G"}' http://localhost:8080/comments -u leo:123
curl -X POST -H "Content-Type: application/json" -d '{"itemId":3, "description": "This is a much, much longer comment by another person mainly for showcasing the possibility of having two comments on one item", "user": "Leonard CarLL"}' http://localhost:8080/comments -u leo:123
# add a comment on item 1
curl -X POST -H "Content-Type: application/json" -d '{"itemId":3, "description": "You know it is coming... This comment will be deleted immediately", "user": "Carl G"}' http://localhost:8080/comments -u leo:123
# delete comment on item 1
curl -X DELETE localhost:8080/comments/3 -u leo:123
