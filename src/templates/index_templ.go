// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/hvossi92/gollama/src/templates/icons"

func Index() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = head().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<body><div class=\"game-container\"><div class=\"game-controls\"><div class=\"control-btn\" title=\"Help\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = icons.Question().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</div><div class=\"control-btn\" title=\"Save Game\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = icons.Save().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</div><a class=\"control-btn\" title=\"Settings\" href=\"/settings\" hx-target=\"body\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = icons.Gear().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</a></div><div class=\"row g-2 mb-2\"><div class=\"col-12\"><!-- Scene Panel --><div class=\"panel scene-panel panel-corner-tl panel-corner-tr panel-corner-bl panel-corner-br\" style=\"background-image: url(&#39;/static/art/baldursgate.png&#39;); \" hx-get=\"/api/scene\" hx-trigger=\"load, sceneChange from:body\" hx-swap=\"outerHTML\"><div class=\"panel-header\">The Forgotten Ruins</div><div class=\"scene-overlay\">Ancient stone columns rise from the mist, their surfaces etched with mysterious runes that glow with a faint blue light.</div></div></div></div><div class=\"row g-2\"><!-- Character Panel - Always visible on mobile, smaller width --><div class=\"col-12 col-md-3 mb-2\" style=\"max-width: 300px;\"><div class=\"panel character-panel panel-corner-tl panel-corner-tr panel-corner-bl panel-corner-br\" hx-get=\"/api/character\" hx-trigger=\"load, characterUpdate from:body\" hx-swap=\"innerHTML\"><div class=\"panel-header\">Character</div><div class=\"character-image\" style=\"background-image: url(&#39;/static/art/tiefling.png&#39;); background-position-y: 0rem;\"></div><div class=\"character-stats\"><h5>Thorne Ironheart</h5><div>Level 5 Dwarf Paladin</div><div class=\"mt-3\">Health</div><div class=\"stat-bar\"><div class=\"stat-fill health-fill\" style=\"width: 75%;\"></div><div class=\"stat-value\">75/100</div></div><div>Mana</div><div class=\"stat-bar\"><div class=\"stat-fill mana-fill\" style=\"width: 60%;\"></div><div class=\"stat-value\">30/50</div></div><div>Stamina</div><div class=\"stat-bar\"><div class=\"stat-fill stamina-fill\" style=\"width: 90%;\"></div><div class=\"stat-value\">45/50</div></div><div class=\"mt-3\"><div class=\"row\"><div class=\"col-6\">STR: 16 (+3)</div><div class=\"col-6\">DEX: 12 (+1)</div><div class=\"col-6\">CON: 15 (+2)</div><div class=\"col-6\">INT: 10 (+0)</div><div class=\"col-6\">WIS: 14 (+2)</div><div class=\"col-6\">CHA: 13 (+1)</div></div></div><div class=\"abilities\"><div class=\"main-btn\" hx-post=\"/api/ability/smite\" hx-trigger=\"click\">Divine Smite</div><div class=\"main-btn\" hx-post=\"/api/ability/heal\" hx-trigger=\"click\">Lay on Hands</div><div class=\"main-btn\" hx-post=\"/api/ability/shield\" hx-trigger=\"click\">Shield of Faith</div></div><div class=\"inventory-section\"><div>Inventory</div><div class=\"inventory-item\" title=\"Health Potion\"><div class=\"item-icon\"><i class=\"fas fa-flask text-danger\"></i></div></div><div class=\"inventory-item\" title=\"Mana Potion\"><div class=\"item-icon\"><i class=\"fas fa-flask text-primary\"></i></div></div><div class=\"inventory-item\" title=\"Warhammer\"><div class=\"item-icon\"><i class=\"fas fa-hammer\"></i></div></div><div class=\"inventory-item\" title=\"Shield\"><div class=\"item-icon\"><i class=\"fas fa-shield-alt\"></i></div></div><div class=\"inventory-item\" title=\"Holy Symbol\"><div class=\"item-icon\"><i class=\"fas fa-sun\"></i></div></div></div><div class=\"effects-section\"><div>Status Effects</div><div class=\"status-effect\" title=\"Blessed\"><div class=\"effect-icon\"><i class=\"fas fa-pray\"></i></div></div><div class=\"status-effect\" title=\"Protected\"><div class=\"effect-icon\"><i class=\"fas fa-shield-alt\"></i></div></div></div></div></div></div><!-- Chat Panel - Full width on mobile --><div class=\"col-12 col-md mb-2\"><div class=\"panel chat-panel panel-corner-tl panel-corner-tr panel-corner-bl panel-corner-br\"><button type=\"button\" hx-post=\"/start\" class=\"main-btn\" style=\"margin: 20%;\" hx-swap=\"outerHTML\" hx-disabled-elt=\"this\" hx-indicator=\"#indicator\">Start your adventure <span id=\"indicator\" class=\"htmx-indicator\">Travelling to the sword coast...</span></button></div></div><!-- Enemy Panel - Full width on mobile --><div class=\"col-12 col-md-3 mb-2\" style=\"max-width: 300px;\"><div class=\"panel enemy-panel panel-corner-tl panel-corner-tr panel-corner-bl panel-corner-br\" hx-get=\"/api/enemy\" hx-trigger=\"load, enemyUpdate from:body\" hx-swap=\"innerHTML\"><div class=\"panel-header\">Enemy</div><div class=\"enemy-image\" style=\"background-image: url(&#39;/static/art/skeleton.png&#39;); background-position-y: 0rem;\"></div><div class=\"enemy-stats\"><h5>Iron Guardian</h5><div>Level 6 Construct</div><div class=\"mt-3\">Health</div><div class=\"stat-bar\"><div class=\"stat-fill health-fill\" style=\"width: 90%;\"></div><div class=\"stat-value\">135/150</div></div><div>Energy</div><div class=\"stat-bar\"><div class=\"stat-fill mana-fill\" style=\"width: 100%;\"></div><div class=\"stat-value\">80/80</div></div><div class=\"mt-3\"><div>Armor Class: 18</div><div>Damage Resistance: Piercing, Slashing</div><div>Weakness: Lightning</div></div><div class=\"mt-3\"><div>Actions:</div><ul class=\"ps-3\"><li>Slam (2d8+4 bludgeoning)</li><li>Stone Projectile (1d10+2 bludgeoning)</li><li>Ground Pound (DC 15 DEX save)</li></ul></div><div class=\"effects-section\"><div>Status Effects</div><div class=\"status-effect\" title=\"Arcane Shield\"><div class=\"effect-icon\"><i class=\"fas fa-shield-alt text-primary\"></i></div></div></div></div></div></div></div></div><script>\n        // Auto-scroll chat to bottom\n        function scrollChatToBottom() {\n            const chatMessages = document.getElementById('chat-messages');\n            if (chatMessages) {\n                chatMessages.scrollTop = chatMessages.scrollHeight;\n            }\n        }\n\n        // Call on page load\n        document.addEventListener('DOMContentLoaded', scrollChatToBottom);\n\n        // Call when new messages are added\n        document.body.addEventListener('htmx:afterSwap', function (event) {\n            if (event.detail.target.id === 'chat-messages') {\n                scrollChatToBottom();\n            }\n        });\n    </script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
