import Script from "next/script";

export default function GroupPage() {
  return (
    <div className="min-h-screen flex flex-col justify-center items-center gap-8">
      <div
        className="tenor-gif-embed"
        data-postid="14878906917560316861"
        data-share-method="host"
        data-aspect-ratio="0.769231"
        data-width="30%"
      >
        <a href="https://tenor.com/view/gorilla-middle-finger-gif-14878906917560316861">
          Gorilla Middle Finger Meme
        </a>
      </div>

      <Script
        src="https://tenor.com/embed.js"
        strategy="afterInteractive"
      />

      <h1 className="text-2xl font-bold">This page is not implemented yet</h1>
    </div>
  );
}