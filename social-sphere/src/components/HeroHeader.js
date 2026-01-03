// 'use client'

// import { useEffect } from 'react'

// export function HeroHeader() {
//   useEffect(() => {
//     // Load Unicorn Studio script only once
//     if (!window.UnicornStudio) {
//       window.UnicornStudio = { isInitialized: false };
//       const script = document.createElement("script");
//       script.src = "https://cdn.jsdelivr.net/gh/hiunicornstudio/unicornstudio.js@v2.0.0/dist/unicornStudio.umd.js";
//       script.onload = function() {
//         if (!window.UnicornStudio.isInitialized) {
//           window.UnicornStudio.init();
//           window.UnicornStudio.isInitialized = true;
//         }
//       };
//       (document.head || document.body).appendChild(script);
//     }
//   }, [])

//   return (
//     <section className="w-full h-[900px] max-h-screen overflow-hidden">
//       <div 
//         data-us-project="MXmygjfLKT9SulOdYdpD" 
//         className="w-full h-full"
//       />
//     </section>
//   )
// }