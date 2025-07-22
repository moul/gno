// This package will be used to manage proof-of-contributions on the exposed smart-contract side.
package rewards

// TODO: write specs.
//
// ROYALTY INTEGRATION:
// This realm should integrate with the royalty fields defined in gnomod.toml files.
// When distributing rewards for package usage or contributions:
// 1. Check if the package has royalties defined in its gnomod.toml
// 2. If royalties exist, distribute rewards according to the weight-based shares
// 3. Apply beneficiary whitelisting to ensure only approved beneficiaries can receive rewards
// 4. If no royalties are defined or beneficiaries are not whitelisted, use a default 
//    alternative treasury address as fallback
//
// The royalty information is stored as:
// - File.Royalties: map[string]int where keys are addresses and values are weights
// - Use GetRoyaltyShares() to convert weights to percentages
// - Example: {"alice": 60, "bob": 40} means alice gets 60% and bob gets 40%
//
// See docs/resources/royalties.md for details.
