"use client";

import { useEffect } from "react";

export default function Modal({ isOpen, onClose, title, description, children, footer }) {
    useEffect(() => {
        if (isOpen) {
            document.body.style.overflow = "hidden";
        } else {
            document.body.style.overflow = "unset";
        }
        return () => {
            document.body.style.overflow = "unset";
        };
    }, [isOpen]);

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm animate-fade-in">
            <div className="bg-(--background) rounded-2xl shadow-2xl w-full max-w-md border border-(--muted)/20 overflow-hidden animate-scale-in">
                <div className="p-6">
                    <div className="flex justify-between items-start mb-4">
                        <h3 className="text-xl font-bold text-(--foreground)">{title}</h3>
                    </div>

                    {description && (
                        <p className="text-(--muted) mb-6 leading-relaxed">
                            {description}
                        </p>
                    )}

                    {children}

                    {footer && (
                        <div className="flex justify-end gap-3 mt-6">
                            {footer}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
