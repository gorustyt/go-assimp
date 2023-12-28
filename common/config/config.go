package config

type Config struct {
	// ---------------------------------------------------------------------------
	/** @brief  Configures whether the AC loader evaluates subdivision surfaces (
	 *  indicated by the presence of the 'subdiv' attribute in the file). By
	 *  default, Assimp performs the subdivision using the standard
	 *  Catmull-Clark algorithm
	 *
	 * * Property type: bool. Default value: true.
	 */
	SubDivision bool
	// ---------------------------------------------------------------------------
	/** @brief  Configures the AC loader to collect all surfaces which have the
	 *    "Backface cull" flag set in separate meshes.
	 *
	 *  Property type: bool. Default value: true.
	 */
	SplitBFCull bool
}

func NewConfig() *Config {
	return &Config{
		SubDivision: true,
		SplitBFCull: true,
	}
}
