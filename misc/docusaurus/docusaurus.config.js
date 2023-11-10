// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Gno.land Documentation",
  favicon: "img/favicon.ico",
  url: "https://docs.gno.land",
  baseUrl: "/",

  organizationName: "gnolang",
  projectName: "gno",

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",

  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          path: "../../docs",
          routeBasePath: "/",
          sidebarPath: require.resolve("./sidebars.js"),
        },
        blog: false,
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        hideOnScroll: true,
        title: "Gno.land",
        logo: {
          alt: "Gno.land Logo",
          src: "img/logo.svg",
          srcDark: "img/logo_light.svg",
        },
        items: [
          {
            type: "docSidebar",
            sidebarId: "tutorialSidebar",
            position: "left",
            label: "Docs",
          },
          {
            href: "https://github.com/gnolang/gno",
            html: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="github-icon">
            <path fill-rule="evenodd" clip-rule="evenodd" d="M12 0.300049C5.4 0.300049 0 5.70005 0 12.3001C0 17.6001 3.4 22.1001 8.2 23.7001C8.8 23.8001 9 23.4001 9 23.1001C9 22.8001 9 22.1001 9 21.1001C5.7 21.8001 5 19.5001 5 19.5001C4.5 18.1001 3.7 17.7001 3.7 17.7001C2.5 17.0001 3.7 17.0001 3.7 17.0001C4.9 17.1001 5.5 18.2001 5.5 18.2001C6.6 20.0001 8.3 19.5001 9 19.2001C9.1 18.4001 9.4 17.9001 9.8 17.6001C7.1 17.3001 4.3 16.3001 4.3 11.7001C4.3 10.4001 4.8 9.30005 5.5 8.50005C5.5 8.10005 5 6.90005 5.7 5.30005C5.7 5.30005 6.7 5.00005 9 6.50005C10 6.20005 11 6.10005 12 6.10005C13 6.10005 14 6.20005 15 6.50005C17.3 4.90005 18.3 5.30005 18.3 5.30005C19 7.00005 18.5 8.20005 18.4 8.50005C19.2 9.30005 19.6 10.4001 19.6 11.7001C19.6 16.3001 16.8 17.3001 14.1 17.6001C14.5 18.0001 14.9 18.7001 14.9 19.8001C14.9 21.4001 14.9 22.7001 14.9 23.1001C14.9 23.4001 15.1 23.8001 15.7 23.7001C20.5 22.1001 23.9 17.6001 23.9 12.3001C24 5.70005 18.6 0.300049 12 0.300049Z" fill="currentColor"/>
            </svg>
            `,
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",

        links: [
          {
            items: [
              {
                html: `
                  <div>
                    <a href="https://gno.land">
                    <svg width="140" height="32" viewBox="0 0 140 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <style>path { fill: var(--ifm-color-primary); }</style>
                    <path fill-rule="evenodd" clip-rule="evenodd" d="M70.1247 0.259277C67.4623 3.35048 63.2383 7.80168 60.4159 10.7585L60.4191 10.7553C54.4159 16.9761 47.7855 23.5041 41.3503 29.2705L40.3583 30.1601L38.6431 31.7409C39.2415 31.1233 40.1087 30.0577 40.9695 28.9505C42.5823 26.8769 43.5775 23.3473 42.3839 21.0049C39.9679 16.2561 40.7391 10.2881 44.7071 6.32328C48.6751 2.35848 54.6399 1.58408 59.3887 4.00008C61.7311 5.19368 65.2607 4.19848 67.3343 2.58568C68.4415 1.72168 69.5071 0.857677 70.1247 0.259277ZM59.6127 7.06248C56.6015 4.85128 52.5567 4.38728 49.1775 6.13128H49.1743C43.1807 9.22888 41.9423 16.4481 45.4623 21.2481C46.0671 21.8337 47.0207 21.8401 47.6319 21.2705L47.6831 21.2193C44.3967 17.9329 44.3967 12.5857 47.6831 9.29928C50.9695 6.01288 56.3167 6.01288 59.6031 9.29928L59.6127 9.28968C60.2271 8.67528 60.2271 7.68008 59.6127 7.06248ZM48.2239 20.6945C48.5823 21.0529 49.1615 21.0529 49.5167 20.6945L49.5199 20.6913C52.6303 17.5617 55.9327 14.2561 59.0655 11.1457C59.4239 10.7873 59.4399 10.1985 59.0815 9.83688L59.0751 9.83048C56.0799 6.83848 51.2095 6.83848 48.2175 9.83048C45.2255 12.8257 45.2255 17.6961 48.2175 20.6881L48.2239 20.6945ZM77.6736 3.29287H75.4656V23.2321H77.6736V3.29287ZM8.85444 8.00644C12.0896 8.00644 14.5952 9.32804 15.9424 11.376V8.54724H18.1536V21.6128C18.1536 26.896 14.192 29.1584 9.44964 29.1584C6.83524 29.1584 3.98084 28.3489 2.39044 27.1905L3.44004 25.4944C4.83844 26.4096 6.80644 27.1905 9.23204 27.1905C13.5456 27.1905 15.9424 25.44 15.9424 21.8816V20.4C14.5952 22.448 11.8176 23.7696 8.85444 23.7696C3.95204 23.7696 0.880035 20.6432 0.880035 15.9008C0.880035 11.1584 3.95204 8.00644 8.85444 8.00644ZM14.4704 20.0704C15.4496 19.1904 15.9424 17.888 15.9424 16.5728H15.9456V15.2288C15.9456 13.7664 15.312 12.3584 14.1664 11.4496C12.8416 10.4 11.2288 9.97444 9.42404 9.97444C5.38244 9.97444 3.12004 12.1312 3.12004 15.904C3.12004 19.6768 5.38244 21.8048 9.42404 21.8048C11.3696 21.8048 13.0944 21.3088 14.4704 20.0704ZM24.3968 11.5904C25.7696 9.37924 28.2496 8.00644 31.4016 8.00644V8.01284C35.9552 8.01284 38.1376 10.9216 38.1376 14.9376V23.2352H35.9296V15.0976C35.9296 12.3776 33.7088 9.98404 30.9888 9.97764C28.7104 9.97124 26.7392 10.8064 25.3504 12.5184C24.72 13.296 24.3968 14.2848 24.3968 15.2896V23.232H22.1888V8.54724H24.3968V11.5904ZM89.6709 8.00646C94.6021 8.00646 97.6197 10.4321 97.6197 15.0945H97.6229V20.3745C97.6229 21.4817 97.8117 22.3713 98.3237 23.2321H95.9781C95.5461 22.5025 95.3573 21.7217 95.3573 20.8865V20.8321C94.1989 22.4225 91.9077 23.7697 88.7013 23.7697C85.4949 23.7697 81.2101 22.4481 81.2101 18.5697C81.2101 14.6913 85.5205 13.3697 88.7013 13.3697C91.8821 13.3697 94.1989 14.6913 95.3573 16.3073V15.0945C95.3573 11.5361 93.3349 9.92007 89.2933 9.92007C86.9477 9.92007 85.1717 10.5665 83.5269 11.6161L82.5029 10.0001C84.2277 8.84167 87.0021 8.00646 89.6709 8.00646ZM88.8901 21.9361C91.3701 21.9361 94.1445 21.1009 95.4117 18.8385V18.8321V18.2945C94.1445 16.0065 91.3701 15.1969 88.8901 15.1969C86.7333 15.1969 83.3925 15.8177 83.3925 18.5665C83.3925 21.3153 86.7333 21.9361 88.8901 21.9361ZM111.414 8.00644C108.262 8.00644 105.782 9.37924 104.409 11.5904V8.54724H102.201V23.2321H104.409V15.2032C104.409 14.3008 104.678 13.408 105.225 12.6912C106.621 10.8608 108.646 9.97125 111.001 9.97445C113.721 9.98085 115.942 12.3744 115.942 15.0944V23.2321H118.15V14.9344C118.15 10.9184 115.968 8.00964 111.414 8.00964V8.00644ZM136.909 11.3761V3.29287H139.12V23.2289H136.909V20.4001C135.562 22.4481 133.053 23.7697 129.821 23.7697C124.918 23.7697 121.846 20.6433 121.846 15.9009C121.846 11.1585 124.918 8.00647 129.821 8.00647C133.056 8.00647 135.562 9.32807 136.909 11.3761ZM136 19.4977C136.592 18.8225 136.909 17.9457 136.909 17.0497V17.0465V14.7233C136.909 13.8273 136.592 12.9505 136 12.2753C134.544 10.6177 132.608 9.97127 130.387 9.97127C126.346 9.97127 124.083 12.1281 124.083 15.9009C124.083 19.6737 126.346 21.8017 130.387 21.8017C132.611 21.8017 134.544 21.1553 136 19.4977ZM69.6511 20.3648C68.7039 20.3648 67.9359 21.1328 67.9359 22.08C67.9359 23.0272 68.7039 23.7952 69.6511 23.7952C70.5983 23.7952 71.3663 23.0272 71.3663 22.08C71.3663 21.1328 70.5983 20.3648 69.6511 20.3648ZM50.2688 22.9281C50.6048 22.6113 51.0816 22.5089 51.5264 22.6369C54.1216 23.3761 57.0336 22.7297 59.0752 20.6881C61.1136 18.6497 61.7632 15.7409 61.024 13.1457C60.896 12.7009 60.9984 12.2241 61.3152 11.8881L61.3504 11.8497C62.7168 14.9409 62.1344 18.6913 59.6064 21.2193C57.0752 23.7505 53.3216 24.3297 50.2304 22.9633L50.2688 22.9281ZM63.1103 13.4977C63.9519 16.6593 63.0463 20.0513 60.7391 22.3585C58.4255 24.6721 55.0207 25.5777 51.8527 24.7201C50.6943 24.4065 49.4911 24.6529 48.5951 25.3729C52.9567 27.9777 58.4703 27.3025 62.0799 23.6961C65.6831 20.0929 66.3615 14.5793 63.7599 10.2177C63.0431 11.1201 62.7999 12.3297 63.1103 13.4977Z"/>
                    </svg>
                    </a>
                  </div>

                  <div class="gno-footer__copy">
                    Made with ♡ by the humans at <a href='https://gno.land'>Gno.land</a>
                  </div>

                  <div class="gno-footer__socials">
                    <a href="https://discord.gg/S8nKUqwkPn">
                      <svg fill="none" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                        <title>Discord</title>
                        <path d="m19.8132 4.87655c2.5253 3.73588 3.7725 7.94985 3.3063 12.80115-.002.0205-.0126.0393-.0294.0517-1.9124 1.4129-3.7652 2.2704-5.592 2.839-.0142.0044-.0295.0041-.0435-.0006-.0141-.0048-.0264-.0139-.035-.0261-.4221-.5908-.8056-1.2138-1.1415-1.8679-.0193-.0385-.0017-.085.038-.1001.609-.2309 1.1881-.5077 1.7452-.8353.0439-.0259.0467-.0891.0061-.1195-.1182-.0883-.2353-.1811-.3474-.2739-.0209-.0172-.0492-.0206-.0729-.009-3.6165 1.6803-7.5782 1.6803-11.23744 0-.02375-.0107-.05198-.0071-.07239.0098-.1118.0928-.22919.1848-.3463.2731-.04053.0304-.03718.0936.00699.1195.55704.3215 1.13617.6044 1.74437.8364.03941.0152.05814.0605.03857.099-.32869.655-.71217 1.2779-1.14205 1.8688-.01872.0239-.04947.0348-.07854.0258-1.81816-.5686-3.67098-1.4261-5.583341-2.839-.015931-.0124-.027391-.0321-.029068-.0526-.389627-4.1962.404439-8.44509 3.303159-12.80109.00699-.01153.01761-.02053.02991-.02587 1.4263-.65865 2.95434-1.14321 4.55142-1.41994.02907-.0045.05814.009.07323.03487.19733.35154.42289.80235.5755 1.17077 1.68348-.25874 3.39318-.25874 5.11178 0 .1527-.36054.3704-.81923.5669-1.17077.007-.01284.0178-.02312.031-.02938.0131-.00626.0279-.00819.0422-.00549 1.5979.27757 3.126.76213 4.5511 1.41994.0126.00534.023.01434.0291.02671zm-9.4762 7.97855c.0176-1.2405-.88132-2.267-2.00967-2.267-1.11913 0-2.00934 1.0175-2.00934 2.267 0 1.2492.90782 2.2667 2.00934 2.2667 1.11941 0 2.00967-1.0175 2.00967-2.2667zm7.4297 0c.0176-1.2405-.8813-2.267-2.0094-2.267-1.1194 0-2.0096 1.0175-2.0096 2.267 0 1.2492.9079 2.2667 2.0096 2.2667 1.1281 0 2.0094-1.0175 2.0094-2.2667z"/>
                      </svg>
                    </a>
                    <a href="https://twitter.com/_gnoland">
                      <svg fill="none" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                        <title>Twitter</title>
                        <path d="M22.0855 5.8404C21.4376 6.11924 20.7564 6.31341 20.0588 6.41811C20.385 6.36221 20.8648 5.77518 21.0559 5.53757C21.3461 5.17916 21.5672 4.77001 21.7081 4.33089C21.7081 4.29828 21.7407 4.25169 21.7081 4.22839C21.6917 4.21942 21.6733 4.21472 21.6545 4.21472C21.6358 4.21472 21.6174 4.21942 21.601 4.22839C20.8435 4.63858 20.0374 4.95165 19.2016 5.16019C19.1725 5.16909 19.1415 5.16989 19.1119 5.1625C19.0824 5.15511 19.0554 5.13982 19.0339 5.11826C18.9688 5.04079 18.8988 4.96765 18.8242 4.89929C18.4833 4.59387 18.0966 4.34389 17.6781 4.15851C17.1132 3.92676 16.5031 3.8264 15.8937 3.86499C15.3025 3.90234 14.7252 4.06092 14.1979 4.33089C13.6785 4.61552 13.2221 5.0022 12.8561 5.46768C12.471 5.94679 12.193 6.5028 12.0407 7.09833C11.9152 7.66479 11.901 8.25025 11.9988 8.82215C11.9988 8.91999 11.9988 8.93396 11.915 8.91999C8.5931 8.43079 5.8676 7.25207 3.64061 4.72225C3.54277 4.61043 3.49152 4.61043 3.41232 4.72225C2.44325 6.19448 2.91381 8.52397 4.12514 9.67474C4.28821 9.82849 4.45593 9.97757 4.63297 10.1173C4.07758 10.0779 3.53576 9.92741 3.0396 9.67474C2.94642 9.61417 2.89517 9.64679 2.89051 9.7586C2.87731 9.91362 2.87731 10.0695 2.89051 10.2245C2.98773 10.9674 3.28051 11.6712 3.73891 12.2638C4.19732 12.8565 4.80491 13.3168 5.49954 13.5976C5.66888 13.6701 5.84532 13.7248 6.02601 13.7607C5.51185 13.8619 4.98452 13.8776 4.46525 13.8073C4.35343 13.784 4.3115 13.8445 4.35343 13.9517C5.0383 15.8153 6.52452 16.3837 7.61472 16.7005C7.7638 16.7238 7.91289 16.7238 8.08062 16.761C8.08062 16.761 8.08061 16.761 8.05266 16.789C7.73119 17.376 6.43134 17.772 5.83499 17.977C4.7465 18.368 3.58597 18.5175 2.43393 18.415C2.25223 18.387 2.2103 18.3917 2.16371 18.415C2.11712 18.4383 2.16371 18.4895 2.21496 18.5361C2.44791 18.6899 2.68086 18.825 2.92313 18.9554C3.64436 19.3488 4.40684 19.6613 5.19671 19.8872C9.28729 21.0147 13.8904 20.1854 16.9606 17.1338C19.374 14.739 20.2219 11.4358 20.2219 8.12796C20.2219 8.00217 20.3757 7.92762 20.4642 7.8624C21.0747 7.38667 21.613 6.82484 22.0622 6.19448C22.14 6.10052 22.1799 5.98088 22.174 5.85904C22.174 5.78915 22.174 5.80313 22.0855 5.8404Z" />
                      </svg>
                    </a>
                    <a href="https://www.youtube.com/@_gnoland">
                      <svg fill="none" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                        <title>YouTube</title>
                        <path d="M18.4309 4H5.56912C3.04566 4 1 6.04566 1 8.56912V14.9986C1 17.5221 3.04566 19.5678 5.56912 19.5678H18.4309C20.9543 19.5678 23 17.5221 23 14.9986V8.56912C23 6.04566 20.9543 4 18.4309 4ZM15.3408 12.0967L9.32495 14.9659C9.16465 15.0424 8.97949 14.9255 8.97949 14.7479V8.83016C8.97949 8.65005 9.16952 8.53333 9.33016 8.61474L15.346 11.6633C15.5249 11.7539 15.5218 12.0104 15.3408 12.0967Z"/>
                      </svg>
                    </a>
                    <a href="https://t.me/gnoland">
                      <svg fill="none" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                        <title>Telegram</title>
                        <path d="M12 0C9.6267 0 7.30654 0.703784 5.33316 2.02236C3.35977 3.34093 1.8217 5.21507 0.913454 7.40779C0.00520429 9.60044 -0.232441 12.0133 0.230579 14.341C0.693599 16.6689 1.83649 18.807 3.51472 20.4853C5.19295 22.1635 7.33114 23.3064 9.65895 23.7694C11.9866 24.2325 14.3995 23.9947 16.5922 23.0865C18.7849 22.1782 20.6591 20.6403 21.9777 18.6669C23.2962 16.6935 24 14.3734 24 12C24 8.8174 22.7357 5.76515 20.4854 3.51472C18.2349 1.26427 15.1826 0 12 0ZM17.895 8.21999L15.93 17.505C15.78 18.165 15.39 18.315 14.835 18.015L11.835 15.795L10.335 17.19C10.2644 17.2822 10.1736 17.3572 10.0697 17.4093C9.9657 17.4612 9.85125 17.4888 9.735 17.49L9.945 14.49L15.495 9.46499C15.75 9.25499 15.495 9.13499 15.135 9.34499L8.325 13.62L5.325 12.69C4.68 12.495 4.665 12.045 5.46 11.745L17.025 7.24499C17.595 7.07999 18.075 7.40999 17.895 8.21999Z"/>
                      </svg>
                    </a>
                  </div>
                  `,
              },
            ],
          },
          {
            title: "Gno Libraries",
            items: [
              {
                label: "gno-js-client",
                href: "https://github.com/gnolang/gno-js-client",
              },
              {
                label: "tm2-js-client",
                href: "https://github.com/gnolang/tm2-js-client",
              },
            ],
          },
          {
            title: "Ecosystem",
            items: [
              {
                label: "Gno by Example",
                href: "https://gno-by-example.com/",
              },
              {
                label: "Awesome GNO",
                href: "https://github.com/gnolang/awesome-gno",
              },
            ],
          },
        ],
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
