import { sendMsg } from "@/actions/chat/send-msg";

export default async function MsgPage() {
    const convs = await sendMsg({interlocutor: "2VolejRejNmG", msg: "ay ay ay ay ayyyyyy puerto ricoooo"});

    if (!convs.success) {
        return (
            <div className="flex flex-col items-center justify-center">
            <h1>Msg not sent</h1>
        </div>
        )
    }

    console.log(convs.data);

    return (
        <div className="flex flex-col items-center justify-center">
            <h1>Hello world</h1>
        </div>
    );
}