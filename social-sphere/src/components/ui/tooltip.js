export default function Tooltip({ content, children }) {
    return (
        <div className="group relative flex items-center justify-center">
            {children}
            <div className="absolute top-full mt-2 px-2.5 py-1 bg-gray-900 text-white text-xs font-medium rounded-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap pointer-events-none z-50 shadow-lg">
                {content}
                {/* Tiny arrow pointing up */}
                <div className="absolute -top-1 left-1/2 -translate-x-1/2 w-2 h-2 bg-gray-900 rotate-45" />
            </div>
        </div>
    );
}
