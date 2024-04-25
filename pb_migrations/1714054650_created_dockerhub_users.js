/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "3j3jf4o9zqakx7j",
    "created": "2024-04-25 14:17:30.009Z",
    "updated": "2024-04-25 14:17:30.009Z",
    "name": "dockerhub_users",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "nins3nkq",
        "name": "name",
        "type": "text",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
        }
      }
    ],
    "indexes": [],
    "listRule": null,
    "viewRule": null,
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("3j3jf4o9zqakx7j");

  return dao.deleteCollection(collection);
})
