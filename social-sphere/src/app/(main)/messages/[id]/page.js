import { getConv } from "@/actions/chat/get-conv";
import { getMessages } from "@/actions/chat/get-messages";
import MessagesContent from "@/components/messages/MessagesContent";

export default async function ConversationPage({ params }) {
    const { id } = await params;

    // Fetch conversations list
    const convsResult = await getConv({ first: true, limit: 50 });
    const conversations = convsResult.success ? convsResult.data : [];

    // Find the selected conversation from the list
    const selectedConversation = conversations.find(
        (conv) => conv.Interlocutor?.id === id
    );

    // Fetch messages for the selected conversation if found
    let initialMessages = [];
    if (selectedConversation) {
        const messagesResult = await getMessages({
            interlocutorId: selectedConversation.Interlocutor?.id,
            limit: 50,
        });
        if (messagesResult.success && messagesResult.data?.Messages) {
            // Messages come newest first, reverse for display
            initialMessages = messagesResult.data.Messages.reverse();
        }
    }

    return (
        <MessagesContent
            initialConversations={conversations}
            initialSelectedId={id}
            initialMessages={initialMessages}
        />
    );
}
