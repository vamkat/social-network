import { getConv } from "@/actions/chat/get-conv";
import MessagesContent from "@/components/messages/MessagesContent";

export default async function MsgPage() {
    const convs = await getConv({ first: true, limit: 15 });

    return (
        <MessagesContent initialConversations={convs.success ? convs.data : []} />
    );
}
