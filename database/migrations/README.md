# Database Migrations

This directory contains SQL migration files for the Kasir API database setup.

## Files

- `001_create_tables.sql` - Creates the `categories` and `products` tables with proper relationships and indexes
- `002_seed_data.sql` - Inserts sample data for categories and products

## Running Migrations on Supabase

### Option 1: Using Supabase Dashboard (Recommended)

1. Go to your **Supabase Project Dashboard**
2. Click **SQL Editor** in the left sidebar
3. Click **New Query**
4. Copy and paste the contents of `001_create_tables.sql`
5. Click **RUN**
6. Create another query and paste `002_seed_data.sql`
7. Click **RUN**

### Option 2: Using psql Command Line

```bash
# For migration file
psql -h aws-1-ap-south-1.pooler.supabase.com \
     -U postgres.ssjxkrawchzwvrqlfzgg \
     -d postgres \
     -f database/migrations/001_create_tables.sql

# For seed data
psql -h aws-1-ap-south-1.pooler.supabase.com \
     -U postgres.ssjxkrawchzwvrqlfzgg \
     -d postgres \
     -f database/migrations/002_seed_data.sql
```

When prompted, enter your database password.

## Schema

### Categories Table
```sql
- id (SERIAL PRIMARY KEY)
- name (VARCHAR 100)
- description (TEXT)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### Products Table
```sql
- id (SERIAL PRIMARY KEY)
- name (VARCHAR 100)
- price (INTEGER)
- stock (INTEGER)
- category_id (INTEGER, FOREIGN KEY)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

## Sample Data

### Categories
1. Makanan - Berbagai jenis makanan instan dan siap saji
2. Minuman - Aneka minuman segar dan kemasan
3. Bumbu Dapur - Perlengkapan bumbu untuk memasak

### Products
1. Indomie Godog - 3500 (Stok: 10) - Makanan
2. Vit 1000ml - 3000 (Stok: 40) - Minuman
3. Kecap - 12000 (Stok: 20) - Bumbu Dapur

## Notes

- The seed data uses `ON CONFLICT ... DO UPDATE` to ensure idempotency (safe to run multiple times)
- Sequences are reset after inserting data to prevent ID conflicts
- Indexes are created for better query performance on `category_id`
