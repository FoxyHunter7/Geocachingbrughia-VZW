<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('message_languages', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->foreignId('message_id')->constrained()->onUpdate('cascade')->onDelete('cascade');
            $table->foreignId('lang_code')->constrained(table: 'languages', column: 'code')->onUpdate('cascade')->onDelete('cascade');
            $table->unique(['message_id','lang_code']);

            $table->string('title');
            $table->string('body')->nullable();

            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('message_languages');
    }
};
