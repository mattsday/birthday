# Data Importer

Download `products.json`:

```bash
wget https://raw.githubusercontent.com/BestBuyAPIs/open-data-set/master/products.json
```

Fix it by replacing any "shipping":"" entries with "shipping":0. You can do this with a simple sed command:

```bash
sed -i 's/"shipping":""/"shipping":0/g' products.json
```

Set your Google Cloud project and collection ID:
```bash
export GOOGLE_CLOUD_PROJECT=<my_project>
export FIRESTORE_COLLECTION=products
```

If not running in Google Cloud, login with the application default credentials:

```bash
gcloud auth application-default login
```

Run the import app:

```bash
go run main.go
```