# Declaration inspection

The declaration inspector reads the stable structural forms of a .gooo source: package, namespace, entity, property, and activity.

    printf '%s\n' '{"source":"package jev\nnamespace example\nentity Receipt\nproperty status string\nactivity observe"}' | go run ./cmd/gooo-declaration-inspect

The result keeps the original declaration digest and adds a shape digest. It rejects orphan properties and incomplete structural entries, while leaving unknown future lines available for forward-compatible parser growth.
