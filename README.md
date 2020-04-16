# Addons for [SQLBoiler](https://github.com/volatiletech/sqlboiler)

## Editor

Generate editor classes to avoid boring column list in inserts/updates (especially if we want to process empty strings, false values, zero numbers).

### Without editors
```
apiKey, err := models.FindAPIKey(ctx, s.db, keyID)
if err != nil {
	...
}
apiKey.Disabled = r.Disabled
apiKey.Title = r.Title
_, err = apiKey.Update(ctx, s.db, boil.Whitelist(models.APIKeyColumns.Disabled, models.APIKeyColumns.Title))
```

### With editors
```
apiKey, err := models.FindAPIKey(ctx, s.db, keyID)
if err != nil {
	...
}
_, err = apiKey.E().SetDisabled(r.Disabled).SetTitle(r.Title).Update(ctx, s.db)
```

### With editors (multiline)
```
apiKey, err := models.FindAPIKey(ctx, s.db, keyID)
if err != nil {
	...
}
_, err = apiKey.E().
	SetDisabled(r.Disabled).
	SetTitle(r.Title).
	Update(ctx, s.db)
```

### Additional details 

The cancelled PR [New "editor" helper to avoid column lists for simple inserts/updates](https://github.com/volatiletech/sqlboiler/pull/629)

### Usage

```
./sqlboiler-addons PATH_TO_MODEL_DIRECTORY
```