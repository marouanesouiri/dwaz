import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Yada',
  description: 'Yet Another Discord API library for Go.',

  locales: {
    root: {
      label: 'English',
      lang: 'en',
      link: '/',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Introduction', link: 'introduction/what-is-yada' }
        ],
        sidebar: [
          {
            text: 'Introduction',
            items: [
              { text: 'What is Yada?', link: 'introduction/what-is-yada' },
              { text: 'Getting Started', link: 'introduction/getting-started' }
            ]
          }
        ]
      }
    },
    fr: {
      label: 'Français',
      lang: 'fr',
      link: '/fr/',
      themeConfig: {
        nav: [
          { text: 'Accueil', link: '/fr/' },
          { text: 'Introduction', link: '/fr/introduction/what-is-yada' }
        ],
        sidebar: [
          {
            text: 'Introduction',
            items: [
              { text: 'Qu\'est-ce que Yada ?', link: '/fr/introduction/what-is-yada' },
              { text: 'Guide de démarrage', link: '/fr/introduction/getting-started' }
            ]
          }
        ]
      }
    }
  }
})
