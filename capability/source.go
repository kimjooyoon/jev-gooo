package capability

import "github.com/kimjooyoon/jev-gooo/envelope"

// DiscoverFromSource binds a .gooo source and performs the same
// declaration-bound discovery as Discover.
func DiscoverFromSource(query, source string) (Discovery, error) {
	declaration, err := envelope.BindDeclaration(source)
	if err != nil {
		return Discovery{}, err
	}
	return Discover(query, declaration)
}
