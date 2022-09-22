import TypeIt from "typeit-react"

export default function TypeItAnimation() {
    return (
             <TypeIt
                getBeforeInit={(instance) => {
                    instance.type("<span class='text-blue-600'>❯</span> ")
                        .type("dockerize").type(" ").pause(750).delete(10).pause(500)
                        .type("terraformize").type(" ").pause(750).delete(13).pause(500)
                        .type("organize").type(" ").pause(750).delete(9).pause(500)
                        .type("standardize").type(" ").pause(750).delete(12).pause(500)
                        .type("optimize").type(" ").pause(750).delete(9).pause(500)
                        .type("ize:").type(" ")

                    return instance;
                 }}                
                options={{
                    speed: 70,
                    waitUntilVisible: true,
                    cursorChar: "|"                    
                }}
                element={"h2"}
                className="text-2xl font-semibold text-gray-800 dark:text-white lg:text-3xl"
                >
                </TypeIt>
    )
}
