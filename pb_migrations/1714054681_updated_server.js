/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("8qgz0s1hmziqchx")

  // remove
  collection.schema.removeField("jli6i8wm")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "zoilclbn",
    "name": "owner",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "3j3jf4o9zqakx7j",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("8qgz0s1hmziqchx")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "jli6i8wm",
    "name": "propietari",
    "type": "relation",
    "required": true,
    "presentable": true,
    "unique": false,
    "options": {
      "collectionId": "_pb_users_auth_",
      "cascadeDelete": true,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  // remove
  collection.schema.removeField("zoilclbn")

  return dao.saveCollection(collection)
})
