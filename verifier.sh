#!/bin/sh
set -eu
cd "$(dirname "$0")"
dossier=$(mktemp -d)
trap 'rm -rf "$dossier"' EXIT

go test ./...
go build -o "$dossier/push-swap" ./cmd/push-swap
go build -o "$dossier/checker" ./cmd/checker

verifier_sorties() {
    code_attendu=$1
    texte_attendu=$2
    erreur_attendue=$3
    shift 3
    code=0
    "$@" > "$dossier/sortie" 2> "$dossier/erreur" || code=$?
    test "$code" -eq "$code_attendu"
    printf '%b' "$texte_attendu" > "$dossier/texte-attendu"
    printf '%b' "$erreur_attendue" > "$dossier/erreur-attendue"
    diff -u "$dossier/texte-attendu" "$dossier/sortie"
    diff -u "$dossier/erreur-attendue" "$dossier/erreur"
}

tri="$dossier/push-swap"
controle="$dossier/checker"
verifier_sorties 0 '' '' "$tri"
verifier_sorties 0 '' '' "$controle" < /dev/null
verifier_sorties 0 '' '' "$tri" '1 2 3'
verifier_sorties 1 '' 'Error\n' "$tri" '2 +02'
verifier_sorties 1 '' 'Error\n' "$controle" '3 un' < /dev/null
printf 'pb\n' > "$dossier/mouvements"
verifier_sorties 0 'KO\n' '' "$controle" '1 2 3' < "$dossier/mouvements"
printf 'sa' > "$dossier/mouvements"
verifier_sorties 1 '' 'Error\n' "$controle" '2 1' < "$dossier/mouvements"

for liste in '4 67 3 87 23' '2 1 3 6 5 8' '2 -3 1 -2 0 -1 3'; do
    "$tri" "$liste" > "$dossier/mouvements"
    verifier_sorties 0 'OK\n' '' "$controle" "$liste" < "$dossier/mouvements"
    printf '%s : %d mouvements, OK\n' "$liste" "$(wc -l < "$dossier/mouvements")"
done
