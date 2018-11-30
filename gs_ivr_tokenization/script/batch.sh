sudo docker rm -f iir-postgres
sudo docker run --name="gss-postgres" -p 5434:5432 -d postgres:9.6-alpine

until sudo docker exec gss-postgres psql -U postgres; do sleep 1; done

until sudo docker exec gss-postgres psql -U banwire -d gs_ivr_tokenization; do sleep 1; sudo docker exec -i iir-postgres psql -U postgres < script/init.sql; done

sudo docker exec -i gss-postgres psql -U banwire -d gs_ivr_tokenization < script/public.sql
