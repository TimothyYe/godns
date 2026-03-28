/** @type {import('next').NextConfig} */
module.exports = () => ({
	output: 'export',
	// Isolate dev artifacts so concurrent `next build` runs do not break `next dev`.
	distDir: process.env.NODE_ENV === 'development' ? '.next-dev' : '.next',
})
