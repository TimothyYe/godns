import nextCoreWebVitals from 'eslint-config-next/core-web-vitals'
import nextTypeScript from 'eslint-config-next/typescript'

const config = [
	{
		ignores: ['.next/**', '.next-dev/**', 'out/**', 'node_modules/**'],
	},
	...nextCoreWebVitals,
	...nextTypeScript,
	{
		rules: {
			'react-hooks/set-state-in-effect': 'off',
		},
	},
]

export default config
