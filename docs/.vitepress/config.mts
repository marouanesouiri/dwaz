import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Dwaz',
  description: 'Discord Wrapper API for Zwafriya for Go.',

  locales: {
    root: {
      label: 'English',
      lang: 'en',
      link: '/',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Introduction', link: 'introduction/what-is-dwaz' }
        ],
        sidebar: [
          {
            text: 'Introduction',
            items: [
              { text: 'What is Dwaz?', link: 'introduction/what-is-dwaz' },
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
          { text: 'Introduction', link: '/fr/introduction/what-is-dwaz' }
        ],
        sidebar: [
          {
            text: 'Introduction',
            items: [
              { text: 'Qu\'est-ce que Dwaz ?', link: '/fr/introduction/what-is-dwaz' },
              { text: 'Guide de démarrage', link: '/fr/introduction/getting-started' }
            ]
          }
        ]
      }
    }
  }
})
