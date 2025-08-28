/** @type {import('tailwindcss').Config} */
module.exports = {
    mode: 'jit',
    content: [
      "./src/**/*.{js,jsx,ts,tsx}",
    ],
    theme: {
      extend: {
        colors: {},
        keyframes: {
          progress: {
            '0%': {
              width: '100%'
            },
            '100%': {
              width: '0%'
            }
          },
          fade: {
            '0%': {
              opacity: 0
            },
            '100%': {
              opacity: 1
            }
          },
          fadeOut: {
            '0%': {
              opacity: 1
            },
            '100%': {
              opacity: 0
            }
          }
        },
        animation: {
          fadeOut: 'fadeOut 0.2s ease-in-out forwards',
          fade: 'fade 0.2s ease-out',
          progress: 'progress 5s linear'
        }
      },
      plugins: [],
    }
  }
  