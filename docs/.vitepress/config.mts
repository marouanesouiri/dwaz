import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Wumpkit',
  description: 'A Lightweight Go library for Discord API.',

  locales: {
    root: {
      label: 'English',
      lang: 'en',
      link: '/',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Introduction', link: 'introduction/what-is-wumpkit' }
        ],
        sidebar: [
          {
            text: 'Introduction',
            items: [
              { text: 'What is Wumpkit?', link: 'introduction/what-is-wumpkit' },
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
          { text: 'Introduction', link: '/fr/introduction/what-is-wumpkit' }
        ],
        sidebar: [
          {
            text: 'Introduction',
            items: [
              { text: 'Qu\'est-ce que Wumpkit ?', link: '/fr/introduction/what-is-wumpkit' },
              { text: 'Guide de démarrage', link: '/fr/introduction/getting-started' }
            ]
          }
        ]
      }
    }
  }
})
