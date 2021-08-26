# online-store

1. Masalah yang terjadi di event 12.12 dikarenakan stock item tidak ditentukan sebagai kolom yang nilainya harus >= 0 sehingga stock tidak mungkin bernilai negatif.
2. Solusi yang saya berikan adalah dengan menentukan stock item harus >= 0. Kemudian setiap user melakukan order baik cart atau checkout harus dilakukan pengecekan stock
   item. Ketika jumlah item yang dipesan melebihi jumlah stock item maka request ditolak. Selain itu, di dalam program yang saya buat stock item hanya berkurang ketika
   user melakukan checkout tidak pada saat cart. Kemudian jika dalam jangka waktu tertentu (dalam program ini saya set 5 menit) user tidak melakukan pembayaran maka status
   order akan diubah menjadi expired dan jumlah stock item yang dipesan akan dikembalikan lagi seperti semula.

# how to run
1. Clone project online-store
  ```bash
  git clone https://github.com/ihsanhusaeri/online-store.git
  ```
2. Install semua package
   ```bash
   go get -u ./...
   ```
3. Install database postgresql
4. Buat database dengan nama **online-store**
5. Isi credential database (db_host, db_user, db_port, db_password) di file **main.go**
6. Jalankan program 
   ```bash
   go run main.go
   ```
# Kekurangan
1. Containerize belum work.
2. Unit testing belum work.
