import { getConv } from "@/actions/chat/get-conv";
import { getMessages } from "@/actions/chat/get-messages";
import MessagesContent from "@/components/messages/MessagesContent";
import { markAsRead } from "@/actions/chat/mark-read";

export default async function ConversationPage({ params }) {
    const { id } = await params;
    let firstMessage = false;

    // Fetch conversations list
    const convsResult = await getConv({ first: true, limit: 15 });
    let conversations = convsResult.success ? convsResult.data : [];

    // Find the selected conversation from the list
    const selectedConversation = conversations.find(
        (conv) => conv.Interlocutor?.id === id
    );

    // Fetch messages for the selected conversation if found
    let initialMessages = [];

    if (selectedConversation) {
        const messagesResult = await getMessages({
            interlocutorId: selectedConversation.Interlocutor?.id,
            limit: 20,
        });
        if (messagesResult.success && messagesResult.data?.Messages) {
            // Messages come newest first, reverse for display
            initialMessages = messagesResult.data.Messages.reverse();
        }

        const res = await markAsRead({convID: selectedConversation.ConversationId, lastMsgID: selectedConversation.LastMessage.id});


    } else {
        firstMessage = true;
    }

    return (
        <MessagesContent
            initialConversations={conversations}
            initialSelectedId={id}
            initialMessages={initialMessages}
            firstMessage={firstMessage}
        />
    );
}
